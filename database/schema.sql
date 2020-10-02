CREATE TABLE parameters (
    parameters_id UUID PRIMARY KEY,
    parameters_name TEXT UNIQUE NOT NULL,
    image_width INTEGER NOT NULL,
    image_height INTEGER NOT NULL,
    file_type TEXT NOT NULL,
    gamma_correction DOUBLE PRECISION NOT NULL,
    texture_gamma DOUBLE PRECISION NOT NULL,
    use_scaling_truncation BOOLEAN NOT NULL,
    samples_per_round INTEGER NOT NULL,
    round_count INTEGER NOT NULL,
    tile_width INTEGER NOT NULL,
    tile_height INTEGER NOT NULL,
    max_bounces INTEGER NOT NULL,
    use_bvh BOOLEAN NOT NULL,
    background_color_magnitude DOUBLE PRECISION NOT NULL,
    background_color DOUBLE PRECISION[3] NOT NULL,
    t_min DOUBLE PRECISION NOT NULL,
    t_max DOUBLE PRECISION NOT NULL
);

CREATE TABLE cameras (
    camera_id UUID PRIMARY KEY,
    camera_name TEXT UNIQUE NOT NULL,
    eye_location DOUBLE PRECISION[3] NOT NULL,
    target_location DOUBLE PRECISION[3] NOT NULL,
    up_vector DOUBLE PRECISION[3] NOT NULL,
    vertical_fov DOUBLE PRECISION NOT NULL,
    aspect_ratio DOUBLE PRECISION NOT NULL,
    aperture DOUBLE PRECISION NOT NULL,
    focus_distance DOUBLE PRECISION NOT NULL
);

CREATE TABLE scenes (
    scene_id UUID PRIMARY KEY,
    scene_name TEXT UNIQUE NOT NULL,
    camera_id TEXT NOT NULL REFERENCES cameras(camera_id)
);

CREATE TYPE OBJECT_TYPE AS ENUM (
    'SPHERE', 
    'CYLINDER', 
    'HOLLOW_CYLINDER', 
    'RECTANGLE',
    'TRIANGLE',
    'PLANE',
    'PYRAMID',
    'BOX',
    'TRANSLATION',
    'ROTATION_X',
    'ROTATION_Y',
    'ROTATION_Z',
    'QUATERNION'
);

CREATE TYPE QUATERNION_ROTATION_ORDER AS ENUM (
    'XYX',
    'XYZ',
    'XZX',
    'XZY',
    'YXY',
    'YXZ',
    'YZY',
    'YZX',
    'ZXY',
    'ZXZ',
    'ZYX',
    'ZYZ'
);

CREATE TABLE objects (
    object_id UUID PRIMARY KEY,
    object_name TEXT UNIQUE NOT NULL,
    object_type OBJECT_TYPE NOT NULL,
    encapsulated_object_id UUID REFERENCES objects(object_id),
    a DOUBLE PRECISION[3],
    b DOUBLE PRECISION[3],
    c DOUBLE PRECISION[3],
    displacement DOUBLE PRECISION[3],
    axis_angles DOUBLE PRECISION[3],
    quaternion_rotation_order QUATERNION_ROTATION_ORDER,
    inner_radius DOUBLE PRECISION,
    outer_radius DOUBLE PRECISION,
    height DOUBLE PRECISION,
    angle DOUBLE PRECISION,
    is_culled BOOLEAN,
    has_negative_normal BOOLEAN,
    has_inverted_normals BOOLEAN,
);

CREATE TYPE TEXTURE_TYPE AS ENUM (
    'COLOR',
    'IMAGE'
);

CREATE TABLE textures (
    texture_id UUID PRIMARY KEY,
    texture_name TEXT UNIQUE NOT NULL,
    color DOUBLE PRECISION[3],
    image_texture BYTEA
);

CREATE TYPE MATERIAL_TYPE AS ENUM (
    'LAMBERTIAN', 
    'METAL', 
    'DIELECTRIC'
);

CREATE TABLE materials (
    material_id UUID PRIMARY KEY,
    material_name TEXT UNIQUE NOT NULL,
    reflectance_texture_id REFERENCES textures(texture_id),
    emittance_texture_id REFERENCES textures(texture_id),
    fuzziness DOUBLE PRECISION,
    refractive_index DOUBLE PRECISION,
);

CREATE TABLE object_materials (
    object_material_id UUID PRIMARY KEY,
    object_id UUID NOT NULL REFERENCES objects(object_id),
    material_id UUID NOT NULL REFERENCES materials(material_id),
    UNIQUE (object_id, material_id)
);

CREATE TABLE scene_object_materials (
    scene_object_material_id UUID PRIMARY KEY,
    scene_id UUID NOT NULL REFERENCES scenes(scene_id),
    object_material_id NOT NULL REFERENCES object_materials(object_material_id),
    UNIQUE (scene_id, object_material_id)
);

CREATE TYPE RENDER_STATUS AS ENUM (
    'CREATED',
    'PENDING',
    'RUNNING',
    'STOPPED',
    'STOPPING',
    'COMPLETED',
    'ERROR'
);

CREATE TABLE renders (
    render_id UUID PRIMARY KEY,
    parameters_id UUID NOT NULL REFERENCES parameters(parameters_id),
    scene_id UUID NOT NULL REFERENCES scenes(scene_id),
    render_status RENDER_STATUS NOT NULL
);

CREATE TYPE FILE_TYPE AS ENUM (
    'PNG',
    'JPEG'
);

CREATE TABLE images (
    image_id UUID PRIMARY KEY,
    render_id UUID NOT NULL REFERENCES renders(render_id),
    file_type FILE_TYPE NOT NULL,
    image_data BYTEA NOT NULL,
    UNIQUE(render_id, file_type)
);