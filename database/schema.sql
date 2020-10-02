CREATE TYPE FILE_TYPE AS ENUM (
    'PNG',
    'JPEG'
);

CREATE TABLE parameters (
    parameters_name TEXT PRIMARY KEY,
    image_width INTEGER NOT NULL,
    image_height INTEGER NOT NULL,
    file_type FILE_TYPE NOT NULL,
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
    camera_name TEXT PRIMARY KEY,
    eye_location DOUBLE PRECISION[3] NOT NULL,
    target_location DOUBLE PRECISION[3] NOT NULL,
    up_vector DOUBLE PRECISION[3] NOT NULL,
    vertical_fov DOUBLE PRECISION NOT NULL,
    aspect_ratio DOUBLE PRECISION NOT NULL,
    aperture DOUBLE PRECISION NOT NULL,
    focus_distance DOUBLE PRECISION NOT NULL
);

CREATE TABLE scenes (
    scene_name TEXT PRIMARY KEY,
    camera_name TEXT NOT NULL REFERENCES cameras(camera_name)
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
    object_name TEXT PRIMARY KEY,
    object_type OBJECT_TYPE NOT NULL,
    encapsulated_object_name TEXT REFERENCES objects(object_name),
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
    has_inverted_normals BOOLEAN
);

CREATE TYPE TEXTURE_TYPE AS ENUM (
    'COLOR',
    'IMAGE'
);

CREATE TABLE textures (
    texture_name TEXT PRIMARY KEY,
    color DOUBLE PRECISION[3],
    image_texture BYTEA
);

CREATE TYPE MATERIAL_TYPE AS ENUM (
    'LAMBERTIAN', 
    'METAL', 
    'DIELECTRIC'
);

CREATE TABLE materials (
    material_name TEXT PRIMARY KEY,
    reflectance_texture_name TEXT REFERENCES textures(texture_name),
    emittance_texture_name TEXT REFERENCES textures(texture_name),
    fuzziness DOUBLE PRECISION,
    refractive_index DOUBLE PRECISION
);

CREATE TABLE object_materials (
    object_name TEXT REFERENCES objects(object_name),
    material_name TEXT REFERENCES materials(material_name),
    PRIMARY KEY (object_name, material_name)
);

CREATE TABLE scene_object_materials (\
    scene_name TEXT REFERENCES scenes(scene_name),
    object_material_name TEXT REFERENCES object_materials(object_material_name),
    PRIMARY KEY (scene_name, object_material_name)
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
    render_name TEXT PRIMARY KEY,
    parameters_name TEXT NOT NULL REFERENCES parameters(parameters_name),
    scene_name TEXT NOT NULL REFERENCES scenes(scene_name),
    render_status RENDER_STATUS NOT NULL,
    image_data BYTEA NOT NULL
);