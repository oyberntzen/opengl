#version 430

layout(local_size_x = 1, local_size_y = 1) in;

layout(rgba32f, binding = 0) uniform image2D raw_texture;

uniform int width;
uniform int height;

uniform float decayRate;
uniform float diffuseRate;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);

    float sum = 0;
    int squares = 0;
    for (int x = -1; x <= 1; x++) {
        for (int y = -1; y <= 1; y++) {
            if (pos.x + x >= 0 && pos.x + x < width && pos.y + y >= 0 && pos.y + y < height) {
                sum += imageLoad(raw_texture, ivec2(pos.x + x, pos.y + y)).rgba.x;
                squares++;
            }
        }
    }
    float strength = (sum * diffuseRate / squares + imageLoad(raw_texture, ivec2(pos.x, pos.y)).rgba.x * (1-diffuseRate)) - decayRate;

    if (strength < 0) {
        strength = 0;
    }

    imageStore(raw_texture, pos, vec4(strength, strength, strength, 1.0));
}