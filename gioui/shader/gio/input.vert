#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

#extension GL_GOOGLE_include_directive : enable

precision highp float;

#include "common.h"

layout(std140,binding=0) uniform Block
{
    // 变换矩阵
    layout(offset=0) vec4 Matrix0;
    layout(offset=16) vec4 Matrix1;
    layout(offset=32) vec4 Matrix2;
    layout(offset=48) vec4 Matrix3;
} _block;

layout(location=0) in vec3 inPos;
layout(location=1) in vec3 inColor;
layout(location=2) in vec2 inUV;
layout(location=3) in vec3 innormal;

layout(location=0) out vec2 vUV;
layout(location=1) out vec3 vertexColor;
layout(location=2) out vec3 vNormal;

void main() {
    vUV = inUV;
    vNormal = innormal; 
    vertexColor = inColor;
    gl_Position =  mat4(_block.Matrix0,_block.Matrix1,_block.Matrix2,_block.Matrix3) * vec4(inPos,1);
}
