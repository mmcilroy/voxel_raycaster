#version 430

const int MAX_RAY_STEPS = 256;

const int RAYCAST_HIT = 0;
const int RAYCAST_MISS = 1;
const int RAYCAST_OUT_OF_BOUNDS = 2;
const int RAYCAST_MAX_STEPS = 3;

const vec4 SKY_BLUE = vec4(102./255., 191./255., 255./255., 1);
const vec4 BROWN = vec4(127./255., 106./255., 79./255., 1);
const vec4 DARK_BROWN = vec4(76./255, 63./255, 47./255, 1);
const vec4 GREEN = vec4(0, 228./255, 48./255, 1);

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Output fragment color
out vec4 finalColor;

// Uniforms
uniform vec2 resolution;
uniform vec3 cameraPos;
uniform vec3 cameraPlane;
uniform vec3 cameraUp;
uniform vec3 cameraRight;
uniform vec3 numVoxels;
uniform vec3 sunPos;

// SSBO
layout (std430, binding=13) buffer voxelData {
    uint voxels[];
};

// Voxel functions
bool isOutside(int x, int y, int z) {
    if (x < 0 || x >= numVoxels.x) {
        return true;
    }
    if (y < 0 || y >= numVoxels.y) {
        return true;
    }
    if (z < 0 || z >= numVoxels.z) {
        return true;
    }
    return false;
}

bool isOutOfBounds(int x, int y, int z, vec3 rayDir) {
	return (x < 0 && rayDir.x <= 0) ||
		(y < 0 && rayDir.y <= 0) ||
		(z < 0 && rayDir.z <= 0) ||
		(x >= numVoxels.x && rayDir.x >= 0) ||
		(y >= numVoxels.y && rayDir.y >= 0) ||
		(z >= numVoxels.z && rayDir.z >= 0);
}

bool getVoxel(int x, int y, int z) {
    // Check if in voxel space
    if (isOutside(x, y, z)) {
        return false;
    }

    // Get the voxel index (which is a byte offset)
    int vi = (x/2) + (z/2)*int(numVoxels.x)/2 + (y/2)*int(numVoxels.x)/2*int(numVoxels.z)/2;

    // Get the corresponding int containing our voxel
    uint vd = voxels[vi/4];

    // Get the byte containing our voxel
    uint vb = vd >> ((vi % 4) * 8) & uint(255);

    // Get the bit containing our voxel
    uint bit = uint(x%2 + (z%2*2) + (y%2*2*2));
    uint mask = uint(1) << bit;

    // Check for voxel presence
    return (vb & mask) != uint(0);
}

int checkHit(int x, int y, int z, vec3 rayDir) {
	if (isOutside(x, y, z)) {
		if (isOutOfBounds(x, y, z, rayDir)) {
			return RAYCAST_OUT_OF_BOUNDS;
		} else {
			return RAYCAST_MISS;
		}
	}

	if (getVoxel(x, y, z)) {
		return RAYCAST_HIT;
	}

	return RAYCAST_MISS;
}

vec3 getRayDir(int x, int y) {
    float aspectRatio = resolution.y / resolution.x;
    vec3 rayPos = cameraPlane + cameraRight * (float(x) / resolution.x - 0.5);
    rayPos = rayPos - cameraUp * ((float(y) / resolution.y - 0.5) * aspectRatio);
    return normalize(rayPos - cameraPos);
}

int raycast(vec3 rayPos, vec3 rayDir, out vec3 hitPos, out ivec3 mapPos) {
    mapPos = ivec3(floor(rayPos + 0.));
    ivec3 rayStep = ivec3(sign(rayDir));
    vec3 deltaDist = abs(vec3(length(rayDir)) / rayDir);
    vec3 sideDist = (sign(rayDir) * (vec3(mapPos) - rayPos) + (sign(rayDir) * 0.5) + 0.5) * deltaDist; 

    int hit = 0;
    int side = 4;
    float dist = 0;

    while (hit == 0) {
        int result = checkHit(mapPos.x, mapPos.y, mapPos.z, rayDir);

		if (result == RAYCAST_OUT_OF_BOUNDS) {
			hit = 0;
            break;
		}

		if (result == RAYCAST_HIT) {
			hit = side;
            break;
		}

        if (sideDist.x < sideDist.y) {
            if (sideDist.x < sideDist.z) {
                dist = sideDist.x;
                sideDist.x += deltaDist.x;
                mapPos.x += rayStep.x;
                side = 1 * rayStep.x;
            }
            else {
                dist = sideDist.z;
                sideDist.z += deltaDist.z;
                mapPos.z += rayStep.z;
                side = 3 * rayStep.z;
            }
        }
        else {
            if (sideDist.y < sideDist.z) {
                dist = sideDist.y;
                sideDist.y += deltaDist.y;
                mapPos.y += rayStep.y;
                side = 2 * rayStep.y;
            }
            else {
                dist = sideDist.z;
                sideDist.z += deltaDist.z;
                mapPos.z += rayStep.z;
                side = 3 * rayStep.z;
            }
        }
    }

    hitPos = rayPos + (rayDir * dist);

	return hit;
}

vec3 hitNormal(int hit) {
	if (hit == -1) {
		return vec3(1, 0, 0);
	} else if (hit == 1) {
		return vec3(-1, 0, 0);
	} else if (hit == -2) {
		return vec3(0, 1, 0);
	} else if (hit == 2) {
		return vec3(0, -1, 0);
	} else if (hit == -3) {
        return vec3(0, 0, 1);
	} else if (hit == 3) {
        return vec3(0, 0, -1);
	} else {
        return vec3(0, 0, 0);
    }
}

float diffuseLight(int hit, vec3 dir) {
	float diffuseLight = dot(hitNormal(hit), dir);
	if (diffuseLight < 0.2) {
		diffuseLight = 0.2;
	}
	return diffuseLight;
}

void main() {
    vec2 fragCoord = fragTexCoord * resolution.xy;
    vec3 rayDir = getRayDir(int(fragCoord.x), int(fragCoord.y));
    vec3 rayPos = cameraPos;
    vec3 hitPos;
    ivec3 mapPos;

    int hit = raycast(rayPos, rayDir, hitPos, mapPos);

    finalColor = SKY_BLUE;
    if (hit == 1 || hit == -1) {
        finalColor = DARK_BROWN;
    } else if (hit == 2 || hit == -2) {
        finalColor = GREEN;
    } else if (hit == 3 || hit == -3) {
        finalColor = BROWN;
    }

    if (hit != 0 && hit != 4) {
        vec3 sunHitPos;
        ivec3 sunMapPos;
        int sunHit = raycast(sunPos, normalize(hitPos - sunPos), sunHitPos, sunMapPos);
        if (sunHit == hit && sunMapPos == mapPos) {
            finalColor = finalColor * diffuseLight(sunHit, normalize(sunPos -sunHitPos));
        } else {
            finalColor = finalColor * 0.2;
        }
    }
}
