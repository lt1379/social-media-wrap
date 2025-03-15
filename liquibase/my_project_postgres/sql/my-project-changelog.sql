 --changeset lamboktulus1379:1 labels:my_project-label context:my_project-context
--comment: my_project.user comment
create table public.user (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(50) not null,
    user_name varchar(50),
    password varchar(50),
    created_by varchar(50),
    updated_by varchar(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
--rollback DROP TABLE public.user;

--changeset lamboktulus1379:2 labels:my_project-label context:my_project-context
--comment: my_project.video comment
 CREATE TABLE IF NOT EXISTS public.video
 (
     id bigserial NOT NULL,
     youtube_video_id character varying(100) NOT NULL,
     created_at timestamp with time zone,
     updated_at timestamp with time zone,
     created_by bigint,
     updated_by bigint,
     youtube_title character varying(100),
     youtube_description text,
     youtube_playlist character varying(100),
     youtube_channel_id character varying(100),
     youtube_channel_username character varying(100),
     youtube_privacy_status character varying(100),
     CONSTRAINT videos_pkey PRIMARY KEY (id)
);
--rollback DROP TABLE public.video;