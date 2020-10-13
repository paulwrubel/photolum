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
    aperture DOUBLE PRECISION NOT NULL,
    focus_distance DOUBLE PRECISION NOT NULL
);

CREATE TABLE scenes (
    scene_name TEXT PRIMARY KEY,
    camera_name TEXT NOT NULL REFERENCES cameras(camera_name)
);

CREATE TYPE PRIMITIVE_TYPE AS ENUM (
    'SPHERE', 
    'CYLINDER', 
    'HOLLOW_CYLINDER', 
    'RECTANGLE',
    'TRIANGLE',
    'PLANE',
    'PYRAMID',
    'BOX',
    'TRANSLATION',
    'ROTATION',
    'QUATERNION',
    'PARTICIPATING_VOLUME'
);

CREATE TYPE AXIS AS ENUM (
    'X',
    'Y',
    'Z'
);

CREATE TYPE ROTATION_ORDER AS ENUM (
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

CREATE TABLE primitives (
    primitive_name TEXT PRIMARY KEY,
    primitive_type PRIMITIVE_TYPE NOT NULL,
    encapsulated_primitive_name TEXT REFERENCES primitives(primitive_name),
    a DOUBLE PRECISION[3],
    b DOUBLE PRECISION[3],
    c DOUBLE PRECISION[3],
    a_normal DOUBLE PRECISION[3],
    b_normal DOUBLE PRECISION[3],
    c_normal DOUBLE PRECISION[3],
    point DOUBLE PRECISION[3],
    normal DOUBLE PRECISION[3],
    center DOUBLE PRECISION[3],
    displacement DOUBLE PRECISION[3],
    axis AXIS,
    axis_angles DOUBLE PRECISION[3],
    rotation_order ROTATION_ORDER,
    radius DOUBLE PRECISION,
    inner_radius DOUBLE PRECISION,
    outer_radius DOUBLE PRECISION,
    height DOUBLE PRECISION,
    angle DOUBLE PRECISION,
    density DOUBLE PRECISION,
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
    texture_type TEXTURE_TYPE NOT NULL,
    color DOUBLE PRECISION[3],
    gamma DOUBLE PRECISION,
    magnitude DOUBLE PRECISION,
    image_data BYTEA
);

CREATE TYPE MATERIAL_TYPE AS ENUM (
    'LAMBERTIAN', 
    'METAL', 
    'DIELECTRIC',
    'ISOTROPIC'
);

CREATE TABLE materials (
    material_name TEXT PRIMARY KEY,
    material_type MATERIAL_TYPE NOT NULL,
    reflectance_texture_name TEXT REFERENCES textures(texture_name),
    emittance_texture_name TEXT REFERENCES textures(texture_name),
    fuzziness DOUBLE PRECISION,
    refractive_index DOUBLE PRECISION,
    CHECK (num_nonnulls(reflectance_texture_name, emittance_texture_name) > 0)
);

CREATE TABLE scene_primitive_materials (
    scene_name TEXT REFERENCES scenes(scene_name),
    primitive_name TEXT REFERENCES primitives(primitive_name),
    material_name TEXT REFERENCES materials(material_name),
    PRIMARY KEY (scene_name, primitive_name, material_name)
);

CREATE TYPE RENDER_STATUS AS ENUM (
    'CREATED',
    'PENDING',
    'STARTING',
    'RUNNING',
    'STOPPING',
    'STOPPED',
    'COMPLETED',
    'ERROR'
);

CREATE TABLE renders (
    render_name TEXT PRIMARY KEY,
    parameters_name TEXT NOT NULL REFERENCES parameters(parameters_name),
    scene_name TEXT NOT NULL REFERENCES scenes(scene_name),
    render_status RENDER_STATUS NOT NULL,
    completed_rounds INTEGER NOT NULL,
    round_progress DOUBLE PRECISION NOT NULL,
    start_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    end_timestamp TIMESTAMP WITH TIME ZONE,
    image_data BYTEA
);