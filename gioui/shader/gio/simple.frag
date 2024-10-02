#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

precision mediump float;

layout(binding = 0) uniform sampler2D vTexture;

layout(location=0) in highp vec2 vUV;
layout(location=1) in highp vec3 vertexColor;
layout(location=2) in highp vec3 vNormal;

layout(location=0) out vec4 fragColor;

void main() {
	vec4 color = texture(vTexture, vec2(vUV));
	vec3 norm = normalize(vNormal);
	float diff = max(dot(norm, vertexColor), 0.2);
	fragColor = vec4(diff * vec3(color[0],color[2],color[1]),color[3]) ;
}
