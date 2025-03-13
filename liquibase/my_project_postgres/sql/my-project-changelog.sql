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

--changeset lamboktulus1379:w labels:my_project-label context:my_project-context
--comment: my_project.video comment
 CREATE TABLE IF NOT EXISTS public.video
 (
     id bigint NOT NULL,
     youtube_video_id bit varying NOT NULL,
     created_at timestamp with time zone,
     updated_at timestamp with time zone,
     created_by bigint,
     updated_by bigint,
     youtube_title bit varying,
     youtube_description bit varying,
     youtube_playlist bit varying,
     youtube_channel_id bit varying,
     youtube_channel_username bit varying,
     youtube_privacy_status bit varying,
     CONSTRAINT videos_pkey PRIMARY KEY (id)
);
--rollback DROP TABLE public.video;