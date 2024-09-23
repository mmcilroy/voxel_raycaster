#version 430

const int WORLD_SIZE = 4; 
const int WORLD_SIZE_CUBED = WORLD_SIZE * WORLD_SIZE * WORLD_SIZE;

in vec4 color;

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// These should be uniforms
vec2 resolution = vec2(1600, 900);
uniform vec3 cameraPos = vec3(0, WORLD_SIZE-1, 0);
uniform float cameraRot = 0;

// SSBO
layout (std430, binding=13) buffer voxelData {
    int voxelsNA[WORLD_SIZE * WORLD_SIZE * WORLD_SIZE];
};

int voxels[WORLD_SIZE * WORLD_SIZE * WORLD_SIZE] = {1,1,1,1,1,1,1,1,0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1,1,0,0,0,0,0};

// Output fragment color
out vec4 finalColor;

const int MAX_RAY_STEPS = 32;

vec2 rotate2d(vec2 v, float a) {
	float sinA = sin(a);
	float cosA = cos(a);
	return vec2(v.x * cosA - v.y * sinA, v.y * cosA + v.x * sinA);	
}

bool isOutside(ivec3 c) {
    if (c.x < 0 || c.x >= WORLD_SIZE ) {
        return true;
    }
    if (c.y < 0 || c.y >= WORLD_SIZE ) {
        return true;
    }
    if (c.z < 0 || c.z >= WORLD_SIZE ) {
        return true;
    }
    return false;
}

bool isOutOfBounds(ivec3 c, vec3 d) {
	return (c.x < 0 && d.x <= 0) ||
		(c.y < 0 && d.y <= 0) ||
		(c.z < 0 && d.z <= 0) ||
		(c.x >= WORLD_SIZE && d.x >= 0) ||
		(c.y >= WORLD_SIZE && d.y >= 0) ||
		(c.z >= WORLD_SIZE && d.z >= 0);
}

bool getVoxel(ivec3 c) {
    if (isOutside(c)) {
        return false;
    }
    int i = c.x + (c.y * WORLD_SIZE) + (c.z * WORLD_SIZE * WORLD_SIZE);
    return i >= 0 && i <= WORLD_SIZE_CUBED-1 && voxels[i] != 0;
}

void main()
{
    vec2 fragCoord = fragTexCoord * resolution.xy;
	vec2 screenPos = -1.0 * ((fragCoord.xy / resolution.xy) * 2.0 - 1.0);
	vec3 cameraDir = vec3(0.0, 0.0, 0.8);
	vec3 cameraPlaneU = vec3(1.0, 0.0, 0.0);
	vec3 cameraPlaneV = vec3(0.0, 1.0, 0.0) * resolution.y / resolution.x;
	vec3 rayDir = cameraDir + screenPos.x * cameraPlaneU + screenPos.y * cameraPlaneV;
    vec3 rayPos = cameraPos;

	rayPos.xz = rotate2d(rayPos.xz, cameraRot);
	rayDir.xz = rotate2d(rayDir.xz, cameraRot);

	ivec3 mapPos = ivec3(floor(rayPos + 0.));

	vec3 deltaDist = abs(vec3(length(rayDir)) / rayDir);
	
	ivec3 rayStep = ivec3(sign(rayDir));

	vec3 sideDist = (sign(rayDir) * (vec3(mapPos) - rayPos) + (sign(rayDir) * 0.5) + 0.5) * deltaDist; 
	 
	bvec3 mask;
	
    bool hit = false;
    
    vec4 color = vec4(0, 0, 0, 1);

	for (int i = 0; i < MAX_RAY_STEPS; i++) {
		if (getVoxel(mapPos)) {
            hit = true;
            break;
        }
        mask = lessThanEqual(sideDist.xyz, min(sideDist.yzx, sideDist.zxy));
        sideDist += vec3(mask) * deltaDist;
        mapPos += ivec3(vec3(mask)) * rayStep;
	}

	finalColor = vec4(0, 0, 0, 1);
    if (hit) {
        //finalColor = vec4(float(mapPos.x)/WORLD_SIZE, float(mapPos.y)/WORLD_SIZE, float(mapPos.z)/WORLD_SIZE, 1);;
        if (mask.x) {
            finalColor = vec4(0.5);
        }
        if (mask.y) {
            finalColor = vec4(1.0);
        }
        if (mask.z) {
            finalColor = vec4(0.75);
        }
    }
}
