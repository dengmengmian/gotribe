--
-- PostgreSQL database dump
--


-- Dumped from database version 16.10 (Debian 16.10-1.pgdg13+1)
-- Dumped by pg_dump version 16.10 (Debian 16.10-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ad; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ad (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(255) NOT NULL,
    description character varying(300) NOT NULL,
    url character varying(255) NOT NULL,
    url_type smallint DEFAULT 1,
    sort smallint DEFAULT 1,
    status smallint DEFAULT 1 NOT NULL,
    scene_id bigint,
    ext text,
    image character varying(255),
    video character varying(255)
);


--
-- Name: COLUMN ad.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.title IS '标题';


--
-- Name: COLUMN ad.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.description IS '描述';


--
-- Name: COLUMN ad.url; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.url IS '广告链接';


--
-- Name: COLUMN ad.url_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.url_type IS '1.链接，2.文章，3.商品';


--
-- Name: COLUMN ad.sort; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.sort IS '排序';


--
-- Name: COLUMN ad.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.status IS '状态，1-未发布；2-发布';


--
-- Name: COLUMN ad.scene_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.scene_id IS '场景 ID';


--
-- Name: COLUMN ad.ext; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.ext IS '扩展字段';


--
-- Name: COLUMN ad.image; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.image IS '图片地址';


--
-- Name: COLUMN ad.video; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad.video IS '视频地址';


--
-- Name: ad_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ad_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ad_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ad_id_seq OWNED BY public.ad.id;


--
-- Name: ad_scene; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ad_scene (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(255) NOT NULL,
    description character varying(300) NOT NULL,
    project_id bigint NOT NULL
);


--
-- Name: COLUMN ad_scene.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad_scene.title IS '标题';


--
-- Name: COLUMN ad_scene.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad_scene.description IS '描述';


--
-- Name: COLUMN ad_scene.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.ad_scene.project_id IS '项目ID';


--
-- Name: ad_scene_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ad_scene_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ad_scene_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ad_scene_id_seq OWNED BY public.ad_scene.id;


--
-- Name: admin; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username character varying(20) NOT NULL,
    password character varying(255) NOT NULL,
    mobile character varying(11) NOT NULL,
    avatar character varying(255),
    nickname character varying(20),
    introduction character varying(255),
    status smallint DEFAULT 1,
    creator character varying(20)
);


--
-- Name: COLUMN admin.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.admin.status IS '1正常, 2禁用';


--
-- Name: admin_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.admin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: admin_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.admin_id_seq OWNED BY public.admin.id;


--
-- Name: admin_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin_roles (
    role_id bigint NOT NULL,
    admin_id bigint NOT NULL
);


--
-- Name: api; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    method character varying(20),
    path character varying(100),
    category character varying(50),
    "desc" character varying(100),
    creator character varying(20)
);


--
-- Name: COLUMN api.method; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.api.method IS '请求方式';


--
-- Name: COLUMN api.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.api.path IS '访问路径';


--
-- Name: COLUMN api.category; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.api.category IS '所属类别';


--
-- Name: COLUMN api."desc"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.api."desc" IS '说明';


--
-- Name: COLUMN api.creator; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.api.creator IS '创建人';


--
-- Name: api_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.api_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: api_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.api_id_seq OWNED BY public.api.id;


--
-- Name: category; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.category (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    parent_id bigint DEFAULT 0,
    sort bigint DEFAULT 1,
    icon character varying(255),
    title character varying(30) NOT NULL,
    slug character varying(30) NOT NULL,
    path character varying(255),
    hidden smallint DEFAULT 1,
    description character varying(300),
    ext text,
    status smallint DEFAULT 1 NOT NULL,
    count bigint DEFAULT 0
);


--
-- Name: COLUMN category.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.parent_id IS '父分类ID，0表示根分类';


--
-- Name: COLUMN category.sort; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.sort IS '排序';


--
-- Name: COLUMN category.icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.icon IS '图标';


--
-- Name: COLUMN category.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.title IS '标题';


--
-- Name: COLUMN category.slug; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.slug IS 'URL别名';


--
-- Name: COLUMN category.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.path IS '自定义url路径';


--
-- Name: COLUMN category.hidden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.hidden IS '1显示，2隐藏';


--
-- Name: COLUMN category.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.description IS '描述';


--
-- Name: COLUMN category.ext; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.ext IS '扩展字段';


--
-- Name: COLUMN category.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.status IS '状态，1-正常；2-禁用';


--
-- Name: COLUMN category.count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.category.count IS '内容数量';


--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.category.id;


--
-- Name: column; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."column" (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    project_id bigint NOT NULL,
    title character varying(30) NOT NULL,
    description character varying(300),
    icon character varying(300),
    info text,
    ext text,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN "column".project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".project_id IS '项目ID';


--
-- Name: COLUMN "column".title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".title IS '标题';


--
-- Name: COLUMN "column".description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".description IS '描述';


--
-- Name: COLUMN "column".icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".icon IS '图片';


--
-- Name: COLUMN "column".info; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".info IS '内容';


--
-- Name: COLUMN "column".ext; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".ext IS '扩展字段';


--
-- Name: COLUMN "column".status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."column".status IS '状态，1-正常；2-禁用';


--
-- Name: column_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.column_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: column_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.column_id_seq OWNED BY public."column".id;


--
-- Name: comment; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.comment (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    project_id bigint NOT NULL,
    content text NOT NULL,
    html_content text NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    object_id character varying(10) NOT NULL,
    object_type smallint DEFAULT 1 NOT NULL,
    type smallint DEFAULT 1 NOT NULL,
    user_id bigint NOT NULL,
    to_user_id bigint NOT NULL,
    parent_id integer DEFAULT 0 NOT NULL,
    reply_to_id integer DEFAULT 0 NOT NULL,
    hot integer DEFAULT 0,
    "like" integer DEFAULT 0,
    dislike integer DEFAULT 0,
    ip character varying(255) NOT NULL,
    country character varying(255) NOT NULL,
    region_name character varying(255) NOT NULL,
    city character varying(255) NOT NULL
);


--
-- Name: COLUMN comment.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.project_id IS '项目ID';


--
-- Name: COLUMN comment.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.content IS '内容';


--
-- Name: COLUMN comment.html_content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.html_content IS 'HTML内容';


--
-- Name: COLUMN comment.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.status IS '状态，1-待审核；2-审核通过';


--
-- Name: COLUMN comment.object_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.object_id IS '评论主题ID';


--
-- Name: COLUMN comment.object_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.object_type IS '评论对象类型，1-文章；2-商品';


--
-- Name: COLUMN comment.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.type IS '评论类型，1-评论；2-回复';


--
-- Name: COLUMN comment.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.user_id IS '用户ID';


--
-- Name: COLUMN comment.to_user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.to_user_id IS '被评论用户ID';


--
-- Name: COLUMN comment.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.parent_id IS '父评论ID';


--
-- Name: COLUMN comment.reply_to_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.reply_to_id IS '回复的评论ID';


--
-- Name: COLUMN comment.hot; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.hot IS '热度';


--
-- Name: COLUMN comment."like"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment."like" IS '点赞数';


--
-- Name: COLUMN comment.dislike; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.dislike IS '踩数';


--
-- Name: COLUMN comment.ip; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.ip IS 'IP地址';


--
-- Name: COLUMN comment.country; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.country IS '国家';


--
-- Name: COLUMN comment.region_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.region_name IS '地区';


--
-- Name: COLUMN comment.city; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comment.city IS '城市';


--
-- Name: comment_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: comment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.comment_id_seq OWNED BY public.comment.id;


--
-- Name: config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.config (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    config_id character varying(10) NOT NULL,
    project_id bigint NOT NULL,
    alias character varying(20) NOT NULL,
    title character varying(30) NOT NULL,
    description character varying(300) NOT NULL,
    type smallint DEFAULT 1 NOT NULL,
    info text NOT NULL,
    md_content text NOT NULL,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN config.config_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.config_id IS '字符ID，分布式 ID';


--
-- Name: COLUMN config.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.project_id IS '项目ID';


--
-- Name: COLUMN config.alias; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.alias IS '别名';


--
-- Name: COLUMN config.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.title IS '标题';


--
-- Name: COLUMN config.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.description IS '描述';


--
-- Name: COLUMN config.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.type IS '类型，1表示普通配置2:json类型';


--
-- Name: COLUMN config.info; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.info IS '内容';


--
-- Name: COLUMN config.md_content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.md_content IS 'MD内容';


--
-- Name: COLUMN config.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.config.status IS '状态，1-正常；2-禁用';


--
-- Name: config_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.config_id_seq OWNED BY public.config.id;


--
-- Name: feedback; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.feedback (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(255) NOT NULL,
    content text,
    phone character varying(20),
    user_id bigint,
    project_id bigint
);


--
-- Name: COLUMN feedback.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.feedback.title IS '标题';


--
-- Name: COLUMN feedback.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.feedback.content IS '内容';


--
-- Name: COLUMN feedback.phone; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.feedback.phone IS '电话';


--
-- Name: COLUMN feedback.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.feedback.user_id IS '用户ID';


--
-- Name: COLUMN feedback.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.feedback.project_id IS '项目 ID';


--
-- Name: feedback_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.feedback_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: feedback_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.feedback_id_seq OWNED BY public.feedback.id;


--
-- Name: menu; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.menu (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(50),
    title character varying(50),
    icon character varying(50),
    path character varying(100),
    redirect character varying(100),
    component character varying(100),
    sort integer DEFAULT 999,
    status smallint DEFAULT 1,
    hidden smallint DEFAULT 2,
    no_cache smallint DEFAULT 2,
    always_show smallint DEFAULT 2,
    breadcrumb smallint DEFAULT 1,
    active_menu character varying(100),
    parent_id bigint DEFAULT 0,
    creator character varying(20)
);


--
-- Name: COLUMN menu.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.name IS '菜单名称(英文名, 可用于国际化)';


--
-- Name: COLUMN menu.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.title IS '菜单标题(无法国际化时使用)';


--
-- Name: COLUMN menu.icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.icon IS '菜单图标';


--
-- Name: COLUMN menu.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.path IS '菜单访问路径';


--
-- Name: COLUMN menu.redirect; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.redirect IS '重定向路径';


--
-- Name: COLUMN menu.component; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.component IS '前端组件路径';


--
-- Name: COLUMN menu.sort; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.sort IS '菜单顺序(1-999)';


--
-- Name: COLUMN menu.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.status IS '菜单状态(正常/禁用, 默认正常)';


--
-- Name: COLUMN menu.hidden; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.hidden IS '菜单在侧边栏隐藏(1隐藏，2显示)';


--
-- Name: COLUMN menu.no_cache; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.no_cache IS '菜单是否被 <keep-alive> 缓存(1不缓存，2缓存)';


--
-- Name: COLUMN menu.always_show; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.always_show IS '忽略之前定义的规则，一直显示根路由(1忽略，2不忽略)';


--
-- Name: COLUMN menu.breadcrumb; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.breadcrumb IS '面包屑可见性(可见/隐藏, 默认可见)';


--
-- Name: COLUMN menu.active_menu; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.active_menu IS '在其它路由时，想在侧边栏高亮的路由';


--
-- Name: COLUMN menu.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.parent_id IS '父菜单编号(编号为0时表示根菜单)';


--
-- Name: COLUMN menu.creator; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.menu.creator IS '创建人';


--
-- Name: menu_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.menu_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: menu_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.menu_id_seq OWNED BY public.menu.id;


--
-- Name: operation_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.operation_log (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username character varying(20),
    ip character varying(45),
    ip_location character varying(128),
    method character varying(20),
    path character varying(100),
    "desc" character varying(100),
    status integer,
    start_time timestamp without time zone,
    time_cost integer,
    user_agent character varying(255)
);


--
-- Name: COLUMN operation_log.username; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.username IS '用户登录名';


--
-- Name: COLUMN operation_log.ip; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.ip IS 'Ip地址';


--
-- Name: COLUMN operation_log.ip_location; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.ip_location IS 'Ip所在地';


--
-- Name: COLUMN operation_log.method; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.method IS '请求方式';


--
-- Name: COLUMN operation_log.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.path IS '访问路径';


--
-- Name: COLUMN operation_log."desc"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log."desc" IS '说明';


--
-- Name: COLUMN operation_log.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.status IS '响应状态码';


--
-- Name: COLUMN operation_log.start_time; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.start_time IS '发起时间';


--
-- Name: COLUMN operation_log.time_cost; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.time_cost IS '请求耗时(ms)';


--
-- Name: COLUMN operation_log.user_agent; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.operation_log.user_agent IS '浏览器标识';


--
-- Name: operation_log_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.operation_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: operation_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.operation_log_id_seq OWNED BY public.operation_log.id;


--
-- Name: point_available; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.point_available (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    project_id bigint NOT NULL,
    user_id bigint,
    points bigint NOT NULL,
    points_log_id bigint NOT NULL,
    expiration_date timestamp with time zone,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN point_available.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.project_id IS '项目ID';


--
-- Name: COLUMN point_available.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.user_id IS '用户ID';


--
-- Name: COLUMN point_available.points; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.points IS '积分数值(分)';


--
-- Name: COLUMN point_available.points_log_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.points_log_id IS '''积分记录表ID''';


--
-- Name: COLUMN point_available.expiration_date; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.expiration_date IS '''过期时间''';


--
-- Name: COLUMN point_available.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_available.status IS '状态，1-正常；2-删除';


--
-- Name: point_available_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.point_available_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: point_available_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.point_available_id_seq OWNED BY public.point_available.id;


--
-- Name: point_deduction; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.point_deduction (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    project_id bigint NOT NULL,
    user_id bigint,
    points bigint NOT NULL,
    points_detail_id bigint,
    available_points_id bigint NOT NULL
);


--
-- Name: COLUMN point_deduction.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_deduction.project_id IS '项目ID';


--
-- Name: COLUMN point_deduction.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_deduction.user_id IS '用户ID';


--
-- Name: COLUMN point_deduction.points; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_deduction.points IS '积分数值(分)';


--
-- Name: COLUMN point_deduction.points_detail_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_deduction.points_detail_id IS '''积分明细ID''';


--
-- Name: COLUMN point_deduction.available_points_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_deduction.available_points_id IS '''可用积分表ID''';


--
-- Name: point_deduction_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.point_deduction_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: point_deduction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.point_deduction_id_seq OWNED BY public.point_deduction.id;


--
-- Name: point_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.point_log (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    project_id bigint NOT NULL,
    user_id bigint,
    points bigint NOT NULL,
    reason character varying(255) NOT NULL,
    type character varying(20) NOT NULL,
    event_id character varying(10),
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN point_log.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.project_id IS '项目ID';


--
-- Name: COLUMN point_log.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.user_id IS '用户ID';


--
-- Name: COLUMN point_log.points; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.points IS '积分数值(分)';


--
-- Name: COLUMN point_log.reason; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.reason IS '加减原因';


--
-- Name: COLUMN point_log.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.type IS '类型';


--
-- Name: COLUMN point_log.event_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.event_id IS '事件ID';


--
-- Name: COLUMN point_log.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.point_log.status IS '状态，1-正常；2-删除';


--
-- Name: point_log_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.point_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: point_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.point_log_id_seq OWNED BY public.point_log.id;


--
-- Name: post; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.post (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    slug character varying(255),
    category_id bigint,
    project_id bigint,
    column_id bigint,
    user_id bigint,
    author character varying(30) NOT NULL,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    html_content text NOT NULL,
    description character varying(300) NOT NULL,
    ext text,
    icon character varying(255),
    view bigint DEFAULT 1,
    type smallint DEFAULT 1,
    is_top smallint DEFAULT 1,
    is_passwd smallint DEFAULT 1,
    pass_word character varying(255) NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    unit_price integer NOT NULL,
    location character varying(255),
    people character varying(255),
    "time" timestamp without time zone,
    images character varying(1000),
    show_time timestamp without time zone,
    video character varying(255) NOT NULL
);


--
-- Name: COLUMN post.slug; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.slug IS 'URL别名/Slug';


--
-- Name: COLUMN post.category_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.category_id IS '分类 ID';


--
-- Name: COLUMN post.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.project_id IS '项目 ID';


--
-- Name: COLUMN post.column_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.column_id IS '专栏ID';


--
-- Name: COLUMN post.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.user_id IS '用户ID';


--
-- Name: COLUMN post.author; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.author IS '作者';


--
-- Name: COLUMN post.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.title IS '标题';


--
-- Name: COLUMN post.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.content IS '内容';


--
-- Name: COLUMN post.html_content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.html_content IS 'html内容';


--
-- Name: COLUMN post.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.description IS '描述';


--
-- Name: COLUMN post.ext; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.ext IS '''扩展字段''';


--
-- Name: COLUMN post.icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.icon IS '图标';


--
-- Name: COLUMN post.view; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.view IS '''阅读量''';


--
-- Name: COLUMN post.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.type IS '类型，1.文章 2.page 3.短文';


--
-- Name: COLUMN post.is_top; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.is_top IS '是否置顶：1-禁用';


--
-- Name: COLUMN post.is_passwd; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.is_passwd IS '是否加密：1-禁用';


--
-- Name: COLUMN post.pass_word; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.pass_word IS '密码';


--
-- Name: COLUMN post.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.status IS '状态，1-草稿；2-发布';


--
-- Name: COLUMN post.unit_price; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.unit_price IS '商品价格(分)';


--
-- Name: COLUMN post.location; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.location IS '地点';


--
-- Name: COLUMN post.people; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.people IS '人物';


--
-- Name: COLUMN post."time"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post."time" IS '业务时间';


--
-- Name: COLUMN post.images; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.images IS '图片';


--
-- Name: COLUMN post.show_time; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.show_time IS '展示时间';


--
-- Name: COLUMN post.video; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post.video IS '产品视频';


--
-- Name: post_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.post_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.post_id_seq OWNED BY public.post.id;


--
-- Name: post_tag; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.post_tag (
    post_id bigint NOT NULL,
    tag_id bigint NOT NULL,
    created_at timestamp with time zone
);


--
-- Name: COLUMN post_tag.post_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post_tag.post_id IS '内容ID';


--
-- Name: COLUMN post_tag.tag_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.post_tag.tag_id IS '标签ID';


--
-- Name: project; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.project (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(30) NOT NULL,
    title character varying(30) NOT NULL,
    description character varying(300),
    keywords character varying(30),
    domain character varying(60),
    post_url character varying(300),
    icp character varying(255),
    public_security character varying(255),
    author character varying(30),
    info text,
    baidu_analytics character varying(255),
    favicon character varying(255),
    nav_image character varying(255),
    push_token character varying(255),
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN project.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.name IS '项目名';


--
-- Name: COLUMN project.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.title IS '网站标题';


--
-- Name: COLUMN project.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.description IS '描述';


--
-- Name: COLUMN project.keywords; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.keywords IS '网站关键词';


--
-- Name: COLUMN project.domain; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.domain IS '项目域名';


--
-- Name: COLUMN project.post_url; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.post_url IS '内容链接';


--
-- Name: COLUMN project.icp; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.icp IS 'icp备案信息';


--
-- Name: COLUMN project.public_security; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.public_security IS '公安备案';


--
-- Name: COLUMN project.author; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.author IS '网站版权';


--
-- Name: COLUMN project.info; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.info IS '内容';


--
-- Name: COLUMN project.baidu_analytics; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.baidu_analytics IS '百度统计';


--
-- Name: COLUMN project.favicon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.favicon IS 'favicon';


--
-- Name: COLUMN project.nav_image; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.nav_image IS '导航图片';


--
-- Name: COLUMN project.push_token; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.push_token IS '百度推送 API token';


--
-- Name: COLUMN project.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.project.status IS '状态，1-正常；2-禁用';


--
-- Name: project_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.project_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: project_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.project_id_seq OWNED BY public.project.id;


--
-- Name: resource; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(255) NOT NULL,
    path character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    file_extension character varying(10) NOT NULL,
    file_type smallint DEFAULT 1 NOT NULL,
    description character varying(300) NOT NULL,
    size integer DEFAULT 0 NOT NULL,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN resource.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.title IS '标题';


--
-- Name: COLUMN resource.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.path IS '路径';


--
-- Name: COLUMN resource.url; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.url IS '当前域名';


--
-- Name: COLUMN resource.file_extension; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.file_extension IS '文件拓展';


--
-- Name: COLUMN resource.file_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.file_type IS '资源类形，1-图片';


--
-- Name: COLUMN resource.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.description IS '描述';


--
-- Name: COLUMN resource.size; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.size IS '文件大小';


--
-- Name: COLUMN resource.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resource.status IS '状态，1-正常；2-禁用';


--
-- Name: resource_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.resource_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: resource_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.resource_id_seq OWNED BY public.resource.id;


--
-- Name: role; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name character varying(20) NOT NULL,
    keyword character varying(20) NOT NULL,
    "desc" character varying(100),
    status smallint DEFAULT 1,
    sort integer DEFAULT 999,
    creator character varying(20)
);


--
-- Name: COLUMN role.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.role.status IS '1正常, 2禁用';


--
-- Name: COLUMN role.sort; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.role.sort IS '角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)';


--
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.role_id_seq OWNED BY public.role.id;


--
-- Name: role_menus; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_menus (
    menu_id bigint NOT NULL,
    role_id bigint NOT NULL
);


--
-- Name: system_config; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.system_config (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    system_config_id character varying(10),
    title character varying(255) NOT NULL,
    content text NOT NULL,
    logo character varying(255),
    icon character varying(255),
    footer character varying(255)
);


--
-- Name: COLUMN system_config.system_config_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.system_config_id IS '唯一字符ID/分布式ID';


--
-- Name: COLUMN system_config.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.title IS '标题';


--
-- Name: COLUMN system_config.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.content IS '内容';


--
-- Name: COLUMN system_config.logo; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.logo IS 'logo';


--
-- Name: COLUMN system_config.icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.icon IS 'icon';


--
-- Name: COLUMN system_config.footer; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_config.footer IS 'footer';


--
-- Name: system_config_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.system_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: system_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.system_config_id_seq OWNED BY public.system_config.id;


--
-- Name: tag; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tag (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title character varying(30) NOT NULL,
    slug character varying(30) NOT NULL,
    description character varying(300),
    color character varying(20),
    sort bigint DEFAULT 1,
    count bigint DEFAULT 0,
    status smallint DEFAULT 1 NOT NULL
);


--
-- Name: COLUMN tag.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.title IS '标题';


--
-- Name: COLUMN tag.slug; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.slug IS 'URL别名';


--
-- Name: COLUMN tag.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.description IS '描述';


--
-- Name: COLUMN tag.color; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.color IS '展示颜色';


--
-- Name: COLUMN tag.sort; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.sort IS '排序，越大越靠前';


--
-- Name: COLUMN tag.count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.count IS '引用次数';


--
-- Name: COLUMN tag.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tag.status IS '状态，1-正常；2-禁用';


--
-- Name: tag_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tag_id_seq OWNED BY public.tag.id;


--
-- Name: third_party_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.third_party_accounts (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint,
    platform character varying(50) NOT NULL,
    bind_flag smallint DEFAULT 1,
    open_id character varying(255) NOT NULL
);


--
-- Name: COLUMN third_party_accounts.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.third_party_accounts.user_id IS '用户ID';


--
-- Name: COLUMN third_party_accounts.platform; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.third_party_accounts.platform IS '平台';


--
-- Name: COLUMN third_party_accounts.bind_flag; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.third_party_accounts.bind_flag IS '是否绑定,2绑定';


--
-- Name: COLUMN third_party_accounts.open_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.third_party_accounts.open_id IS 'openID';


--
-- Name: third_party_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.third_party_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: third_party_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.third_party_accounts_id_seq OWNED BY public.third_party_accounts.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."user" (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username character varying(30) NOT NULL,
    project_id bigint NOT NULL,
    password character varying(255) NOT NULL,
    nickname character varying(30) NOT NULL,
    email character varying(254) DEFAULT NULL::character varying,
    phone character varying(32) DEFAULT NULL::character varying,
    sex character(1) DEFAULT 'M'::bpchar NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    birthday date,
    background character varying(255),
    ext text,
    avatar_url character varying(255)
);


--
-- Name: COLUMN "user".username; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".username IS '用户名';


--
-- Name: COLUMN "user".project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".project_id IS '项目ID';


--
-- Name: COLUMN "user".password; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".password IS '密码';


--
-- Name: COLUMN "user".nickname; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".nickname IS '昵称';


--
-- Name: COLUMN "user".email; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".email IS '邮箱';


--
-- Name: COLUMN "user".phone; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".phone IS '电话';


--
-- Name: COLUMN "user".sex; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".sex IS 'M:男 F:女';


--
-- Name: COLUMN "user".status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".status IS '用户状态，1-正常；2-禁用';


--
-- Name: COLUMN "user".birthday; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".birthday IS '用户生日，格式为YYYY-MM-DD';


--
-- Name: COLUMN "user".background; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".background IS '个人中心背景';


--
-- Name: COLUMN "user".ext; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".ext IS '扩展字段';


--
-- Name: COLUMN "user".avatar_url; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."user".avatar_url IS '头像地址';


--
-- Name: user_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_event (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint NOT NULL,
    project_id bigint NOT NULL,
    event_type smallint DEFAULT 1 NOT NULL,
    event_detail text,
    duration integer,
    ip character varying(255),
    user_agent character varying(255),
    referer character varying(255),
    platform character varying(255)
);


--
-- Name: COLUMN user_event.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.user_id IS '用户ID';


--
-- Name: COLUMN user_event.project_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.project_id IS '项目ID';


--
-- Name: COLUMN user_event.event_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.event_type IS '事件类型，1-浏览事件；2-点击事件';


--
-- Name: COLUMN user_event.event_detail; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.event_detail IS '事件详情';


--
-- Name: COLUMN user_event.duration; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.duration IS '事件时长';


--
-- Name: COLUMN user_event.ip; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.ip IS 'IP地址';


--
-- Name: COLUMN user_event.user_agent; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.user_agent IS '用户代理';


--
-- Name: COLUMN user_event.referer; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.referer IS '来源页面';


--
-- Name: COLUMN user_event.platform; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_event.platform IS '平台';


--
-- Name: user_event_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_event_id_seq OWNED BY public.user_event.id;


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: ad id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ad ALTER COLUMN id SET DEFAULT nextval('public.ad_id_seq'::regclass);


--
-- Name: ad_scene id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ad_scene ALTER COLUMN id SET DEFAULT nextval('public.ad_scene_id_seq'::regclass);


--
-- Name: admin id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin ALTER COLUMN id SET DEFAULT nextval('public.admin_id_seq'::regclass);


--
-- Name: api id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api ALTER COLUMN id SET DEFAULT nextval('public.api_id_seq'::regclass);


--
-- Name: category id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);


--
-- Name: column id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."column" ALTER COLUMN id SET DEFAULT nextval('public.column_id_seq'::regclass);


--
-- Name: comment id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comment ALTER COLUMN id SET DEFAULT nextval('public.comment_id_seq'::regclass);


--
-- Name: config id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config ALTER COLUMN id SET DEFAULT nextval('public.config_id_seq'::regclass);


--
-- Name: feedback id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.feedback ALTER COLUMN id SET DEFAULT nextval('public.feedback_id_seq'::regclass);


--
-- Name: menu id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu ALTER COLUMN id SET DEFAULT nextval('public.menu_id_seq'::regclass);


--
-- Name: operation_log id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operation_log ALTER COLUMN id SET DEFAULT nextval('public.operation_log_id_seq'::regclass);


--
-- Name: point_available id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_available ALTER COLUMN id SET DEFAULT nextval('public.point_available_id_seq'::regclass);


--
-- Name: point_deduction id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_deduction ALTER COLUMN id SET DEFAULT nextval('public.point_deduction_id_seq'::regclass);


--
-- Name: point_log id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_log ALTER COLUMN id SET DEFAULT nextval('public.point_log_id_seq'::regclass);


--
-- Name: post id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post ALTER COLUMN id SET DEFAULT nextval('public.post_id_seq'::regclass);


--
-- Name: project id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project ALTER COLUMN id SET DEFAULT nextval('public.project_id_seq'::regclass);


--
-- Name: resource id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource ALTER COLUMN id SET DEFAULT nextval('public.resource_id_seq'::regclass);


--
-- Name: role id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role ALTER COLUMN id SET DEFAULT nextval('public.role_id_seq'::regclass);


--
-- Name: system_config id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_config ALTER COLUMN id SET DEFAULT nextval('public.system_config_id_seq'::regclass);


--
-- Name: tag id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag ALTER COLUMN id SET DEFAULT nextval('public.tag_id_seq'::regclass);


--
-- Name: third_party_accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.third_party_accounts ALTER COLUMN id SET DEFAULT nextval('public.third_party_accounts_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Name: user_event id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_event ALTER COLUMN id SET DEFAULT nextval('public.user_event_id_seq'::regclass);


--
-- Name: ad ad_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ad
    ADD CONSTRAINT ad_pkey PRIMARY KEY (id);


--
-- Name: ad_scene ad_scene_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ad_scene
    ADD CONSTRAINT ad_scene_pkey PRIMARY KEY (id);


--
-- Name: admin admin_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (id);


--
-- Name: admin_roles admin_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_roles
    ADD CONSTRAINT admin_roles_pkey PRIMARY KEY (role_id, admin_id);


--
-- Name: api api_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api
    ADD CONSTRAINT api_pkey PRIMARY KEY (id);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- Name: column column_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."column"
    ADD CONSTRAINT column_pkey PRIMARY KEY (id);


--
-- Name: comment comment_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comment
    ADD CONSTRAINT comment_pkey PRIMARY KEY (id);


--
-- Name: config config_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.config
    ADD CONSTRAINT config_pkey PRIMARY KEY (id);


--
-- Name: feedback feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_pkey PRIMARY KEY (id);


--
-- Name: menu menu_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu
    ADD CONSTRAINT menu_pkey PRIMARY KEY (id);


--
-- Name: operation_log operation_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.operation_log
    ADD CONSTRAINT operation_log_pkey PRIMARY KEY (id);


--
-- Name: point_available point_available_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_available
    ADD CONSTRAINT point_available_pkey PRIMARY KEY (id);


--
-- Name: point_deduction point_deduction_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_deduction
    ADD CONSTRAINT point_deduction_pkey PRIMARY KEY (id);


--
-- Name: point_log point_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.point_log
    ADD CONSTRAINT point_log_pkey PRIMARY KEY (id);


--
-- Name: post post_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post
    ADD CONSTRAINT post_pkey PRIMARY KEY (id);


--
-- Name: post_tag post_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_tag
    ADD CONSTRAINT post_tag_pkey PRIMARY KEY (post_id, tag_id);


--
-- Name: project project_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.project
    ADD CONSTRAINT project_pkey PRIMARY KEY (id);


--
-- Name: resource resource_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource
    ADD CONSTRAINT resource_pkey PRIMARY KEY (id);


--
-- Name: role_menus role_menus_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_menus
    ADD CONSTRAINT role_menus_pkey PRIMARY KEY (menu_id, role_id);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_pkey PRIMARY KEY (id);


--
-- Name: system_config system_config_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_config
    ADD CONSTRAINT system_config_pkey PRIMARY KEY (id);


--
-- Name: tag tag_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_pkey PRIMARY KEY (id);


--
-- Name: third_party_accounts third_party_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.third_party_accounts
    ADD CONSTRAINT third_party_accounts_pkey PRIMARY KEY (id);


--
-- Name: admin uni_admin_mobile; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT uni_admin_mobile UNIQUE (mobile);


--
-- Name: admin uni_admin_username; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT uni_admin_username UNIQUE (username);


--
-- Name: role uni_role_keyword; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT uni_role_keyword UNIQUE (keyword);


--
-- Name: role uni_role_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT uni_role_name UNIQUE (name);


--
-- Name: user_event user_event_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_event
    ADD CONSTRAINT user_event_pkey PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: idx_ad_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ad_deleted_at ON public.ad USING btree (deleted_at);


--
-- Name: idx_ad_scene_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ad_scene_deleted_at ON public.ad_scene USING btree (deleted_at);


--
-- Name: idx_ad_scene_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ad_scene_id ON public.ad USING btree (scene_id);


--
-- Name: idx_ad_scene_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ad_scene_project_id ON public.ad_scene USING btree (project_id);


--
-- Name: idx_admin_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_admin_deleted_at ON public.admin USING btree (deleted_at);


--
-- Name: idx_api_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_deleted_at ON public.api USING btree (deleted_at);


--
-- Name: idx_api_path_method; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_api_path_method ON public.api USING btree (path, method);


--
-- Name: idx_category_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_category_deleted_at ON public.category USING btree (deleted_at);


--
-- Name: idx_category_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_category_parent_id ON public.category USING btree (parent_id);


--
-- Name: idx_category_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_category_slug ON public.category USING btree (slug);


--
-- Name: idx_category_sort; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_category_sort ON public.category USING btree (sort);


--
-- Name: idx_column_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_column_deleted_at ON public."column" USING btree (deleted_at);


--
-- Name: idx_column_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_column_project_id ON public."column" USING btree (project_id);


--
-- Name: idx_comment_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_deleted_at ON public.comment USING btree (deleted_at);


--
-- Name: idx_comment_object_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_object_id ON public.comment USING btree (object_id);


--
-- Name: idx_comment_object_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_object_type ON public.comment USING btree (object_type);


--
-- Name: idx_comment_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_project_id ON public.comment USING btree (project_id);


--
-- Name: idx_comment_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_status ON public.comment USING btree (status);


--
-- Name: idx_comment_to_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_to_user_id ON public.comment USING btree (to_user_id);


--
-- Name: idx_comment_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comment_user_id ON public.comment USING btree (user_id);


--
-- Name: idx_config_alias; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_config_alias ON public.config USING btree (alias);


--
-- Name: idx_config_config_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_config_config_id ON public.config USING btree (config_id);


--
-- Name: idx_config_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_config_deleted_at ON public.config USING btree (deleted_at);


--
-- Name: idx_config_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_config_project_id ON public.config USING btree (project_id);


--
-- Name: idx_feedback_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedback_deleted_at ON public.feedback USING btree (deleted_at);


--
-- Name: idx_feedback_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedback_project_id ON public.feedback USING btree (project_id);


--
-- Name: idx_feedback_title; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_feedback_title ON public.feedback USING btree (title);


--
-- Name: idx_feedback_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedback_user_id ON public.feedback USING btree (user_id);


--
-- Name: idx_menu_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_menu_deleted_at ON public.menu USING btree (deleted_at);


--
-- Name: idx_operation_log_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operation_log_deleted_at ON public.operation_log USING btree (deleted_at);


--
-- Name: idx_operation_log_ip_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operation_log_ip_time ON public.operation_log USING btree (ip, start_time DESC);


--
-- Name: idx_operation_log_path_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operation_log_path_time ON public.operation_log USING btree (path, start_time DESC);


--
-- Name: idx_operation_log_status_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operation_log_status_time ON public.operation_log USING btree (status, start_time DESC);


--
-- Name: idx_operation_log_username_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_operation_log_username_time ON public.operation_log USING btree (username, start_time DESC);


--
-- Name: idx_point_available_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_available_deleted_at ON public.point_available USING btree (deleted_at);


--
-- Name: idx_point_available_project_user_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_available_project_user_status ON public.point_available USING btree (project_id, user_id, status);


--
-- Name: idx_point_available_user_status_expiration; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_available_user_status_expiration ON public.point_available USING btree (user_id, status, expiration_date);


--
-- Name: idx_point_deduction_available_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_deduction_available_user ON public.point_deduction USING btree (available_points_id, user_id);


--
-- Name: idx_point_deduction_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_deduction_deleted_at ON public.point_deduction USING btree (deleted_at);


--
-- Name: idx_point_deduction_project_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_deduction_project_user_created_at ON public.point_deduction USING btree (project_id, user_id, created_at);


--
-- Name: idx_point_log_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_log_deleted_at ON public.point_log USING btree (deleted_at);


--
-- Name: idx_point_log_project_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_log_project_created_at ON public.point_log USING btree (project_id, created_at);


--
-- Name: idx_point_log_project_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_log_project_user_created_at ON public.point_log USING btree (project_id, user_id, created_at);


--
-- Name: idx_point_log_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_point_log_user_created_at ON public.point_log USING btree (user_id, created_at);


--
-- Name: idx_post_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_category_id ON public.post USING btree (category_id);


--
-- Name: idx_post_column_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_column_id ON public.post USING btree (column_id);


--
-- Name: idx_post_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_deleted_at ON public.post USING btree (deleted_at);


--
-- Name: idx_post_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_project_id ON public.post USING btree (project_id);


--
-- Name: idx_post_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_post_slug ON public.post USING btree (slug);


--
-- Name: idx_post_tag_tag_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_tag_tag_id ON public.post_tag USING btree (tag_id);


--
-- Name: idx_post_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_user_id ON public.post USING btree (user_id);


--
-- Name: idx_project_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_project_deleted_at ON public.project USING btree (deleted_at);


--
-- Name: idx_resource_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resource_deleted_at ON public.resource USING btree (deleted_at);


--
-- Name: idx_role_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_deleted_at ON public.role USING btree (deleted_at);


--
-- Name: idx_system_config_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_system_config_deleted_at ON public.system_config USING btree (deleted_at);


--
-- Name: idx_system_config_system_config_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_system_config_system_config_id ON public.system_config USING btree (system_config_id);


--
-- Name: idx_system_config_title; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_system_config_title ON public.system_config USING btree (title);


--
-- Name: idx_tag_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tag_deleted_at ON public.tag USING btree (deleted_at);


--
-- Name: idx_tag_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tag_slug ON public.tag USING btree (slug);


--
-- Name: idx_tag_title; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_tag_title ON public.tag USING btree (title);


--
-- Name: idx_third_party_accounts_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_third_party_accounts_deleted_at ON public.third_party_accounts USING btree (deleted_at);


--
-- Name: idx_third_party_accounts_open_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_third_party_accounts_open_id ON public.third_party_accounts USING btree (open_id);


--
-- Name: idx_third_party_accounts_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_third_party_accounts_user_id ON public.third_party_accounts USING btree (user_id);


--
-- Name: idx_user_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_deleted_at ON public."user" USING btree (deleted_at);


--
-- Name: idx_user_event_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_event_deleted_at ON public.user_event USING btree (deleted_at);


--
-- Name: idx_user_event_event_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_event_event_type ON public.user_event USING btree (event_type);


--
-- Name: idx_user_event_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_event_project_id ON public.user_event USING btree (project_id);


--
-- Name: idx_user_project_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_user_project_email ON public."user" USING btree (project_id, email);


--
-- Name: idx_user_project_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_project_id ON public."user" USING btree (project_id);


--
-- Name: idx_user_project_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_user_project_phone ON public."user" USING btree (project_id, phone);


--
-- Name: idx_user_project_username; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_user_project_username ON public."user" USING btree (project_id, username);


--
-- Name: idx_username; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_username ON public.post USING btree (author);


--
-- PostgreSQL database dump complete
--


