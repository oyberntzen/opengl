#version 430

#define pi 3.14159
#define numberOfAgents 50000

layout(local_size_x = 1) in;
layout(rgba32f, binding = 0) uniform image2D raw_texture;
layout(std430, binding = 1) buffer agents 
{
    float xPositions[numberOfAgents];
    float yPositions[numberOfAgents];
    float angles[numberOfAgents];
};

uniform int width;
uniform int height;

uniform uint time;
uniform float trailWeight;
uniform float moveSpeed;
uniform float turnSpeed;
uniform float sensorAngle;
uniform float sensorDistance;
uniform int sensorSize;

uint hash( uint x ) {
    x += ( x << 10u );
    x ^= ( x >>  6u );
    x += ( x <<  3u );
    x ^= ( x >> 11u );
    x += ( x << 15u );
    return x;
}



// Compound versions of the hashing algorithm I whipped together.
uint hash( uvec2 v ) { return hash( v.x ^ hash(v.y)                         ); }
uint hash( uvec3 v ) { return hash( v.x ^ hash(v.y) ^ hash(v.z)             ); }
uint hash( uvec4 v ) { return hash( v.x ^ hash(v.y) ^ hash(v.z) ^ hash(v.w) ); }



// Construct a float with half-open range [0:1] using low 23 bits.
// All zeroes yields 0.0, all ones yields the next smallest representable value below 1.0.
float floatConstruct( uint m ) {
    const uint ieeeMantissa = 0x007FFFFFu; // binary32 mantissa bitmask
    const uint ieeeOne      = 0x3F800000u; // 1.0 in IEEE binary32

    m &= ieeeMantissa;                     // Keep only mantissa bits (fractional part)
    m |= ieeeOne;                          // Add fractional part to 1.0

    float  f = uintBitsToFloat( m );       // Range [1:2]
    return f - 1.0;                        // Range [0:1]
}



// Pseudo-random value in half-open range [0:1].
float random( float x ) { return floatConstruct(hash(floatBitsToUint(x))); }
float random( uint  x ) { return floatConstruct(hash(x));                  }
float random( ivec2 v ) { return floatConstruct(hash(floatBitsToUint(v))); }
float random( ivec3 v ) { return floatConstruct(hash(floatBitsToUint(v))); }
float random( ivec4 v ) { return floatConstruct(hash(floatBitsToUint(v))); }

float sense(vec2 pos, float angle) {
    vec2 direction = vec2(cos(angle), sin(angle)) * sensorDistance;
    vec2 sensorPos = ivec2(int((pos.x + direction.x)*width), int((pos.y + direction.y)*height));

    float sum = 0;
    int squares = 0;
    for (int x = -sensorSize; x <= sensorSize; x++) {
        for (int y = -sensorSize; y <= sensorSize; y++) {
            if (sensorPos.x + x >= 0 && sensorPos.x + x < width && sensorPos.y + y >= 0 && sensorPos.y + y < height) {
                sum += imageLoad(raw_texture, ivec2(sensorPos.x + x, sensorPos.y + y)).rgba.x;
                squares++;
            }
        }
    }
    return sum / squares;
}

void main() {
    int index = int(gl_GlobalInvocationID.x);

    float xPosition = xPositions[index];
    float yPosition = yPositions[index];
    float angle = angles[index];

    float rand = random(time * index);

    vec2 pos = vec2(xPosition, yPosition);
    float weightForward = sense(pos, angle);
    float weightLeft = sense(pos, angle+sensorAngle);
    float weightRight = sense(pos, angle-sensorAngle);

    rand = random(rand);
    float randomSteer = rand;

    if (weightForward > weightLeft && weightForward > weightRight) {
        angle += 0;
    }
    else if (weightForward < weightLeft && weightForward < weightRight) {
        angle += (randomSteer - 0.5) * 2 * turnSpeed;
    }
    else if (weightRight > weightLeft) {
        angle -= randomSteer * turnSpeed;
    }
    else if (weightLeft > weightRight) {
        angle += randomSteer * turnSpeed;
    }

    float xNew = xPosition + cos(angle) * moveSpeed;
    float yNew = yPosition + sin(angle) * moveSpeed;

    if (xNew < 0) {
        xNew = 0.0;
        float newAngle = (random(rand) - 0.5) * pi;
        angles[index] = newAngle;
    }
    if (xNew > 1) {
        xNew = 1.0;
        float newAngle = (random(rand) + 0.5) * pi;
        angles[index] = newAngle;
    }
    if (yNew < 0) {
        yNew = 0.0;
        float newAngle = (random(rand)) * pi;
        angles[index] = newAngle;
    }
    if (yNew > 1) {
        yNew = 1.0;
        float newAngle = (random(rand) + 1.0) * pi;
        angles[index] = newAngle;
    }

    xPositions[index] = xNew;
    yPositions[index] = yNew;

    ivec2 position = ivec2(int(xNew*float(width)), int(yNew*float(height)));
    
    float strength = imageLoad(raw_texture, position).x + trailWeight;
    if (strength > 1.0) {
        strength = 1;
    }

    imageStore(raw_texture, position, vec4(strength, strength, strength, 1.0));
}