CREATE TYPE masa_kerja_type AS (
	r double precision,
	i double precision
);

CREATE FUNCTION calc_age(param_tanggal date) RETURNS smallint
    LANGUAGE plpgsql
    AS $$BEGIN
	--Routine body goes here...

	RETURN (
			(
				date_part('year' :: TEXT,(now()) :: DATE) - date_part(
					'year' :: TEXT,
					param_tanggal
				)
			) * (12) :: DOUBLE PRECISION
		) + (
			date_part(
				'month' :: TEXT,
				(now()) :: DATE
			) - date_part(
				'month' :: TEXT,
				param_tanggal
			)
		);
END
$$;

CREATE FUNCTION get_kgb_yad(__tmt_cpns date) RETURNS date
    LANGUAGE plpgsql
    AS $$
	declare dt_yad date;
	declare y_tmt int;
	declare m_tmt int;
	declare d_tmt int;
	declare y_yad int;
	declare m_yad int;
	declare d_yad int;
	BEGIN
	    y_tmt = date_part('year',__tmt_cpns);
			m_tmt = date_part('month',__tmt_cpns);
			d_tmt = date_part('day',__tmt_cpns);

			y_yad	= date_part('year',CURRENT_DATE);

			dt_yad = (y_yad||'-'||m_tmt||'-'||d_tmt)::date;
			if(mod(y_yad,2) = mod(y_tmt,2)) then -- sama -genap
				if(dt_yad >=current_date) then
						return dt_yad;
				else
						return (dt_yad + interval '2 year')::date;
				end if;
				else
						y_yad = y_yad +1;
						dt_yad = (y_yad||'-'||m_tmt||'-'||d_tmt)::date;
						return dt_yad;
			end if;
	RETURN dt_yad;
END
$$;

CREATE FUNCTION get_masa_kerja(date_start date, date_end date) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
	declare _month varchar;
	BEGIN
	-- Routine body goes here...

	select  extract(year from age(date_end,date_start)) || ' thn ' || extract(month from age(date_end,date_start)) || ' bln 'into _month;
	return _month;
END
$$;

CREATE FUNCTION get_masa_kerja_arr(date_start date, date_end date) RETURNS integer[]
    LANGUAGE plpgsql
    AS $$
	declare _month int;
	declare _year int;
	declare _json int[2];
	BEGIN
		select extract(year from age(date_end,date_start)) ,extract(month from age(date_end,date_start))
		into _year,_month;


	 _json[0]= _year;
	 _json[1]= _month;
	--select json_build_object('year'::text,extract(year from age(date_end,date_start))::int,'month'::text,extract(month from age(date_end,date_start))) into _json;
	-- select json_build_object('year',1,'month',2) into _json;
	return _json;
END
$$;

CREATE FUNCTION get_month_masa_kerja(date_start date, date_end date) RETURNS smallint
    LANGUAGE plpgsql
    AS $$
	declare _month int4;
	BEGIN
	-- Routine body goes here...

	select  extract(year from age(date_end,date_start))*12 + extract(month from age(date_end,date_start)) into _month;
	return _month;
END
$$;

CREATE FUNCTION uuid_generate_v1() RETURNS uuid
    LANGUAGE c STRICT
    AS '$libdir/uuid-ossp', 'uuid_generate_v1';

CREATE FUNCTION uuid_generate_v1mc() RETURNS uuid
    LANGUAGE c STRICT
    AS '$libdir/uuid-ossp', 'uuid_generate_v1mc';

CREATE FUNCTION uuid_generate_v3(namespace uuid, name text) RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_generate_v3';

CREATE FUNCTION uuid_generate_v4() RETURNS uuid
    LANGUAGE c STRICT
    AS '$libdir/uuid-ossp', 'uuid_generate_v4';

CREATE FUNCTION uuid_generate_v5(namespace uuid, name text) RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_generate_v5';

CREATE FUNCTION uuid_nil() RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_nil';

CREATE FUNCTION uuid_ns_dns() RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_ns_dns';

CREATE FUNCTION uuid_ns_oid() RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_ns_oid';

CREATE FUNCTION uuid_ns_url() RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_ns_url';

CREATE FUNCTION uuid_ns_x500() RETURNS uuid
    LANGUAGE c IMMUTABLE STRICT
    AS '$libdir/uuid-ossp', 'uuid_ns_x500';

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE rwt_nine_box (
    "ID" integer NOT NULL,
    "PNS_NIP" character varying(32),
    "NAMA" character varying(200),
    "NAMA_JABATAN" character varying(200),
    "KELAS_JABATAN" smallint,
    "KESIMPULAN" character varying(255),
    "TAHUN" character varying(4)
);

CREATE SEQUENCE "NINE_BOX_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "NINE_BOX_ID_seq" OWNED BY rwt_nine_box."ID";

CREATE TABLE absen (
    "ID" integer NOT NULL,
    "NIP" character varying(30),
    "NAMA" character varying(200),
    "TANGGAL" date,
    "JAM" character varying(20),
    "VERIFIKASI" smallint DEFAULT 0,
    latitude character varying(30),
    longitude character varying(30),
    inside_office smallint,
    input_type smallint,
    keterangan character varying(255),
    is_wfo smallint,
    timezoned character varying(100),
    waktu timestamp(6) without time zone,
    updated_at timestamp(6) without time zone,
    created_at timestamp(6) without time zone
);

COMMENT ON COLUMN absen."VERIFIKASI" IS '1=veririkasi ok, 0 : blm verifikasi';

COMMENT ON COLUMN absen.input_type IS '2 untuk pake browser';

CREATE SEQUENCE "absen_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "absen_ID_seq" OWNED BY absen."ID";

CREATE SEQUENCE activities_activity_id_seq
    START WITH 52450
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE activities (
    activity_id integer DEFAULT nextval('activities_activity_id_seq'::regclass) NOT NULL,
    user_id bigint DEFAULT (0)::bigint NOT NULL,
    activity text NOT NULL,
    module character varying(255) NOT NULL,
    created_on timestamp(0) without time zone,
    deleted integer DEFAULT 0 NOT NULL
);

CREATE SEQUENCE agama_id_seq
    START WITH 7
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE agama (
    "ID" smallint DEFAULT nextval('agama_id_seq'::regclass) NOT NULL,
    "NAMA" character varying(20),
    "NCSISTIME" character varying(30),
    deleted smallint
);

CREATE TABLE anak (
    "ID" bigint NOT NULL,
    "PASANGAN" bigint,
    "NAMA" character varying(255),
    "JENIS_KELAMIN" character varying(1),
    "TANGGAL_LAHIR" date,
    "TEMPAT_LAHIR" character varying(255),
    "STATUS_ANAK" character varying(1),
    "PNS_ID" character varying(32),
    "NIP" character varying(30)
);

CREATE SEQUENCE "anak_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "anak_ID_seq" OWNED BY anak."ID";

CREATE TABLE arsip (
    "ID" integer NOT NULL,
    "ID_JENIS_ARSIP" integer,
    "NIP" character varying(25),
    "KETERANGAN" character varying(255),
    "EXTENSION_FILE" character varying(100),
    "JENIS_FILE" character varying(100),
    "FILE_SIZE" character varying(20),
    "FILE_BASE64" text,
    "CREATED_BY" integer,
    "CREATED_DATE" date,
    "UPDATED_BY" integer,
    "UPDATED_DATE" date,
    "ISVALID" integer DEFAULT 0,
    location character varying(255),
    name character varying(255),
    sk_number character varying(100),
    ref uuid DEFAULT uuid_generate_v4()
);

COMMENT ON COLUMN arsip."EXTENSION_FILE" IS '.doc, .xls, dll';

COMMENT ON COLUMN arsip."JENIS_FILE" IS 'image, document, zip, pdf';

CREATE SEQUENCE "arsip_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "arsip_ID_seq" OWNED BY arsip."ID";

CREATE TABLE asesmen_hasil_asesmen (
    id bigint NOT NULL,
    nip character varying,
    jenis_asesmen_jabatan character varying,
    satuan_kerja character varying,
    tanggal_asesmen date,
    jpm numeric,
    integritas numeric,
    kerjasama numeric,
    komunikasi numeric,
    orientasi_pada_hasil numeric,
    pelayanan_publik numeric,
    pengembangan_diri_dan_orang_lain numeric,
    mengelola_perubahan numeric,
    pengambilan_keputusan numeric,
    perekat_bangsa numeric
);

CREATE SEQUENCE asesmen_hasil_asesmen_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE asesmen_hasil_asesmen_id_seq OWNED BY asesmen_hasil_asesmen.id;

CREATE TABLE asesmen_pegawai_berpotensi_jpt (
    id bigint NOT NULL,
    nip character varying,
    usia character varying,
    status_kepegawaian character varying,
    golongan character varying,
    jenis_jabatan character varying,
    jabatan character varying,
    tmt character varying,
    lama_jabatan_terakhir character varying,
    eselon character varying,
    satker character varying,
    unit_organisasi_induk character varying,
    kedudukan character varying,
    tipe character varying,
    pendidikan character varying,
    jabatan_madya_lain character varying,
    tmt_jabatan_madya_lain character varying,
    jabatan_struktural_lain character varying,
    tmt_jabatan_struktural_lain character varying,
    lama_menjabat_akumulasi bigint,
    rekam_jejak text,
    skp text,
    asesmen character varying,
    hukuman_disiplin character varying,
    jabatan_struktural_lainnya_json text,
    nama character varying
);

CREATE SEQUENCE asesmen_pegawai_berpotensi_jpt_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE asesmen_pegawai_berpotensi_jpt_id_seq OWNED BY asesmen_pegawai_berpotensi_jpt.id;

CREATE TABLE asesmen_riwayat_hukuman_disiplin (
    id bigint NOT NULL,
    nip character varying,
    tingkat_hukuman_disiplin character varying,
    jenis_hukuman_disiplin character varying,
    no_sk character varying,
    tanggal_sk character varying,
    status character varying,
    tahun integer,
    alasan text
);

CREATE SEQUENCE asesmen_riwayat_hukuman_disiplin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE asesmen_riwayat_hukuman_disiplin_id_seq OWNED BY asesmen_riwayat_hukuman_disiplin.id;

CREATE TABLE baperjakat (
    "ID" integer NOT NULL,
    "TANGGAL" date NOT NULL,
    "KETERANGAN" character varying(50),
    "TANGGAL_PENETAPAN" date,
    "NO_SK_PENETAPAN" character varying(20),
    "STATUS_AKTIF" integer,
    "TANGGAL_PELANTIKAN" date
);

CREATE SEQUENCE "baperjakat_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "baperjakat_ID_seq" OWNED BY baperjakat."ID";

CREATE SEQUENCE pegawai_id_seq
    START WITH 32781
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE pegawai (
    "ID" integer DEFAULT nextval('pegawai_id_seq'::regclass) NOT NULL,
    "PNS_ID" character varying(36) NOT NULL,
    "NIP_LAMA" character varying(9),
    "NIP_BARU" character varying(18),
    "NAMA" character varying(255),
    "GELAR_DEPAN" character varying(60),
    "GELAR_BELAKANG" character varying(60),
    "TEMPAT_LAHIR_ID" character varying(50),
    "TGL_LAHIR" date,
    "JENIS_KELAMIN" character varying(10),
    "AGAMA_ID" integer,
    "JENIS_KAWIN_ID" character varying(255),
    "NIK" character varying(255),
    "NOMOR_DARURAT" character varying(255),
    "NOMOR_HP" character varying(60),
    "EMAIL" character varying(255),
    "ALAMAT" character varying(255),
    "NPWP" character varying(255),
    "BPJS" character varying(50),
    "JENIS_PEGAWAI_ID" character varying(50),
    "KEDUDUKAN_HUKUM_ID" character varying(36),
    "STATUS_CPNS_PNS" character varying(20),
    "KARTU_PEGAWAI" character varying(30),
    "NOMOR_SK_CPNS" character varying(60),
    "TGL_SK_CPNS" date,
    "TMT_CPNS" date,
    "TMT_PNS" date,
    "GOL_AWAL_ID" character varying(36),
    "GOL_ID" integer,
    "TMT_GOLONGAN" date,
    "MK_TAHUN" character varying(20),
    "MK_BULAN" character varying(20),
    "JENIS_JABATAN_IDx" character varying(36),
    "JABATAN_ID" character varying(36),
    "TMT_JABATAN" date,
    "PENDIDIKAN_ID" character varying(36),
    "TAHUN_LULUS" character varying(20),
    "KPKN_ID" character varying(36),
    "LOKASI_KERJA_ID" character varying(36),
    "UNOR_ID" character varying(36),
    "UNOR_INDUK_ID" character varying(36),
    "INSTANSI_INDUK_ID" character varying(36),
    "INSTANSI_KERJA_ID" character varying(36),
    "SATUAN_KERJA_INDUK_ID" character varying(36),
    "SATUAN_KERJA_KERJA_ID" character varying(36),
    "GOLONGAN_DARAH" character varying(20),
    "PHOTO" character varying(100),
    "TMT_PENSIUN" date,
    "LOKASI_KERJA" character(200),
    "JML_ISTRI" character(1),
    "JML_ANAK" character(1),
    "NO_SURAT_DOKTER" character(100),
    "TGL_SURAT_DOKTER" date,
    "NO_BEBAS_NARKOBA" character(100),
    "TGL_BEBAS_NARKOBA" date,
    "NO_CATATAN_POLISI" character(100),
    "TGL_CATATAN_POLISI" date,
    "AKTE_KELAHIRAN" character(50),
    "STATUS_HIDUP" character(15),
    "AKTE_MENINGGAL" character(50),
    "TGL_MENINGGAL" date,
    "NO_ASKES" character(50),
    "NO_TASPEN" character(50),
    "TGL_NPWP" date,
    "TEMPAT_LAHIR" character(200),
    "PENDIDIKAN" character(165),
    "TK_PENDIDIKAN" character(3),
    "TEMPAT_LAHIR_NAMA" character(200),
    "JENIS_JABATAN_NAMA" character(200),
    "JABATAN_NAMA" character(254),
    "KPKN_NAMA" character(255),
    "INSTANSI_INDUK_NAMA" character(100),
    "INSTANSI_KERJA_NAMA" character(160),
    "SATUAN_KERJA_INDUK_NAMA" character(170),
    "SATUAN_KERJA_NAMA" character(155),
    "JABATAN_INSTANSI_ID" character(15),
    "BUP" smallint DEFAULT 58,
    "JABATAN_INSTANSI_NAMA" character varying(512) DEFAULT NULL::character varying,
    "JENIS_JABATAN_ID" integer,
    terminated_date date,
    status_pegawai smallint DEFAULT 1,
    "JABATAN_PPNPN" character(255),
    "JABATAN_INSTANSI_REAL_ID" character(36),
    "CREATED_DATE" date,
    "CREATED_BY" integer,
    "UPDATED_DATE" date,
    "UPDATED_BY" integer,
    "EMAIL_DIKBUD_BAK" character varying(255),
    "EMAIL_DIKBUD" character varying(100),
    "KODECEPAT" character varying(100),
    "IS_DOSEN" smallint,
    "MK_TAHUN_SWASTA" smallint DEFAULT 0,
    "MK_BULAN_SWASTA" smallint DEFAULT 0,
    "KK" character varying(30),
    "NIDN" character varying(30),
    "KET" character varying(255),
    "NO_SK_PEMBERHENTIAN" character varying(100),
    status_pegawai_backup smallint,
    "MASA_KERJA" character varying,
    "KARTU_ASN" character varying
);

COMMENT ON COLUMN pegawai.status_pegawai IS '1=pns,2=honorer';

CREATE TABLE unitkerja (
    "NO" character varying(255),
    "KODE_INTERNAL" character varying(255),
    "ID" character varying(255) NOT NULL,
    "NAMA_UNOR" character varying(255),
    "ESELON_ID" character varying(255),
    "CEPAT_KODE" character varying(255),
    "NAMA_JABATAN" character varying(255),
    "NAMA_PEJABAT" character varying(255),
    "DIATASAN_ID" character varying(255),
    "INSTANSI_ID" character varying(255),
    "PEMIMPIN_NON_PNS_ID" character varying(255),
    "PEMIMPIN_PNS_ID" character varying(255),
    "JENIS_UNOR_ID" character varying(255),
    "UNOR_INDUK" character varying(255),
    "JUMLAH_IDEAL_STAFF" character varying(255),
    "ORDER" bigint,
    deleted smallint,
    "IS_SATKER" smallint DEFAULT 0 NOT NULL,
    "ESELON_1" character varying(32),
    "ESELON_2" character varying(32),
    "ESELON_3" character varying(32),
    "ESELON_4" character varying(32),
    "EXPIRED_DATE" date,
    "KETERANGAN" character varying(255),
    "JENIS_SATKER" character varying(255),
    "ABBREVIATION" character varying(255),
    "UNOR_INDUK_PENYETARAAN" character varying(255),
    "JABATAN_ID" character varying(32),
    "WAKTU" character varying(4),
    "PERATURAN" character varying(100)
);

CREATE MATERIALIZED VIEW vw_unit_list AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(32) AS "ESELON_1",
    (x.eselon[2])::character varying(32) AS "ESELON_2",
    (x.eselon[3])::character varying(32) AS "ESELON_3",
    (x.eselon[4])::character varying(32) AS "ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WHERE (uk."EXPIRED_DATE" IS NULL)
  WITH NO DATA;

CREATE VIEW daftar_pegawai AS
 SELECT pegawai."ID",
    pegawai."PNS_ID",
    pegawai."NIP_LAMA",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."KEDUDUKAN_HUKUM_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."KARTU_PEGAWAI",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."GOL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."JENIS_JABATAN_IDx" AS "JENIS_JABATAN_ID",
    pegawai."JABATAN_ID",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."TAHUN_LULUS",
    pegawai."KPKN_ID",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."PHOTO",
    pegawai."TMT_PENSIUN",
    pegawai."LOKASI_KERJA",
    pegawai."NO_SURAT_DOKTER",
    pegawai."JML_ISTRI",
    pegawai."JML_ANAK",
    pegawai."TGL_SURAT_DOKTER",
    pegawai."NO_BEBAS_NARKOBA",
    pegawai."TGL_BEBAS_NARKOBA",
    pegawai."NO_CATATAN_POLISI",
    pegawai."TGL_CATATAN_POLISI",
    pegawai."AKTE_KELAHIRAN",
    pegawai."STATUS_HIDUP",
    pegawai."AKTE_MENINGGAL",
    pegawai."TGL_MENINGGAL",
    pegawai."NO_ASKES",
    pegawai."NO_TASPEN",
    pegawai."TGL_NPWP",
    pegawai."TEMPAT_LAHIR",
    pegawai."PENDIDIKAN",
    pegawai."TK_PENDIDIKAN",
    pegawai."TEMPAT_LAHIR_NAMA",
    pegawai."JENIS_JABATAN_NAMA",
    pegawai."JABATAN_NAMA",
    pegawai."KPKN_NAMA",
    pegawai."INSTANSI_INDUK_NAMA",
    pegawai."INSTANSI_KERJA_NAMA",
    pegawai."SATUAN_KERJA_INDUK_NAMA",
    pegawai."SATUAN_KERJA_NAMA",
    pegawai."JABATAN_INSTANSI_ID",
    pegawai."BUP",
    date_part('year'::text, age((pegawai."TGL_LAHIR")::timestamp with time zone)) AS "AGE",
    vw."ESELON_1",
    vw."ESELON_2",
    vw."ESELON_3",
    vw."ESELON_4",
    date_part('year'::text, pegawai."TGL_LAHIR") AS tahun_lahir
   FROM (pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)));

CREATE VIEW daftar_pns_aktif AS
 SELECT "ID",
    "PNS_ID",
    "NIP_LAMA",
    "NIP_BARU",
    "NAMA",
    "GELAR_DEPAN",
    "GELAR_BELAKANG",
    "TEMPAT_LAHIR_ID",
    "TGL_LAHIR",
    "JENIS_KELAMIN",
    "AGAMA_ID",
    "JENIS_KAWIN_ID",
    "NIK",
    "NOMOR_DARURAT",
    "NOMOR_HP",
    "EMAIL",
    "ALAMAT",
    "NPWP",
    "BPJS",
    "JENIS_PEGAWAI_ID",
    "KEDUDUKAN_HUKUM_ID",
    "STATUS_CPNS_PNS",
    "KARTU_PEGAWAI",
    "NOMOR_SK_CPNS",
    "TGL_SK_CPNS",
    "TMT_CPNS",
    "TMT_PNS",
    "GOL_AWAL_ID",
    "GOL_ID",
    "TMT_GOLONGAN",
    "MK_TAHUN",
    "MK_BULAN",
    "JENIS_JABATAN_IDx" AS "JENIS_JABATAN_ID",
    "JABATAN_ID",
    "TMT_JABATAN",
    "PENDIDIKAN_ID",
    "TAHUN_LULUS",
    "KPKN_ID",
    "LOKASI_KERJA_ID",
    "UNOR_ID",
    "UNOR_INDUK_ID",
    "INSTANSI_INDUK_ID",
    "INSTANSI_KERJA_ID",
    "SATUAN_KERJA_INDUK_ID",
    "SATUAN_KERJA_KERJA_ID",
    "GOLONGAN_DARAH",
    "PHOTO",
    "TMT_PENSIUN",
    "LOKASI_KERJA",
    "NO_SURAT_DOKTER",
    "JML_ISTRI",
    "JML_ANAK",
    "TGL_SURAT_DOKTER",
    "NO_BEBAS_NARKOBA",
    "TGL_BEBAS_NARKOBA",
    "NO_CATATAN_POLISI",
    "TGL_CATATAN_POLISI",
    "AKTE_KELAHIRAN",
    "STATUS_HIDUP",
    "AKTE_MENINGGAL",
    "TGL_MENINGGAL",
    "NO_ASKES",
    "NO_TASPEN",
    "TGL_NPWP",
    "TEMPAT_LAHIR",
    "PENDIDIKAN",
    "TK_PENDIDIKAN",
    "TEMPAT_LAHIR_NAMA",
    "JENIS_JABATAN_NAMA",
    "JABATAN_NAMA",
    "KPKN_NAMA",
    "INSTANSI_INDUK_NAMA",
    "INSTANSI_KERJA_NAMA",
    "SATUAN_KERJA_INDUK_NAMA",
    "SATUAN_KERJA_NAMA",
    "JABATAN_INSTANSI_ID",
    "BUP",
    date_part('year'::text, age(("TGL_LAHIR")::timestamp with time zone)) AS "AGE"
   FROM pegawai pegawai
  WHERE ((status_pegawai = 1) AND ((terminated_date IS NULL) OR ((terminated_date IS NOT NULL) AND (terminated_date > ('now'::text)::date))));

CREATE TABLE daftar_rohaniawan (
    id integer NOT NULL,
    nip character varying(30),
    nama character varying(100),
    jabatan character varying(100),
    agama integer,
    aktif character varying(5),
    pangkat_gol character varying(30)
);

CREATE SEQUENCE daftar_rohaniawan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE daftar_rohaniawan_id_seq OWNED BY daftar_rohaniawan.id;

CREATE TABLE golongan (
    "NAMA" character varying(255),
    "NAMA_PANGKAT" character varying(255),
    "ID" integer NOT NULL,
    "NAMA2" character varying(255),
    "GOL" smallint,
    "GOL_PPPK" character varying(255)
);

CREATE SEQUENCE "jabatan_No_seq"
    START WITH 1776
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE jabatan (
    "NO" integer DEFAULT nextval('"jabatan_No_seq"'::regclass) NOT NULL,
    "KODE_JABATAN" character varying NOT NULL,
    "NAMA_JABATAN" character varying,
    "NAMA_JABATAN_FULL" text,
    "JENIS_JABATAN" character varying,
    "KELAS" smallint,
    "PENSIUN" smallint,
    "KODE_BKN" character varying(32),
    id bigint NOT NULL,
    "NAMA_JABATAN_BKN" character varying(255),
    "KATEGORI_JABATAN" character varying(100),
    "BKN_ID" character varying(36)
);

CREATE TABLE kedudukan_hukum (
    "ID" character varying(4) NOT NULL,
    "NAMA" character varying(255)
);

CREATE VIEW data_pegawai AS
 SELECT a."NIP_BARU",
    a."NAMA",
    a."JABATAN_NAMA",
    b."KELAS",
    c."NAMA" AS gol,
    c."NAMA_PANGKAT",
    d."NAMA_UNOR",
    a."STATUS_CPNS_PNS",
    e."NAMA" AS status,
    f."NAMA_UNOR" AS satker
   FROM (((((pegawai a
     JOIN jabatan b ON ((a."JABATAN_INSTANSI_ID" = (b."KODE_JABATAN")::bpchar)))
     JOIN golongan c ON ((a."GOL_ID" = c."ID")))
     JOIN unitkerja d ON (((a."UNOR_ID")::text = (d."ID")::text)))
     JOIN unitkerja f ON (((f."ID")::text = (d."UNOR_INDUK")::text)))
     JOIN kedudukan_hukum e ON (((a."KEDUDUKAN_HUKUM_ID")::text = (e."ID")::text)))
  WHERE (a."JABATAN_INSTANSI_ID" <> '0'::bpchar);

CREATE TABLE hari_libur (
    "ID" integer NOT NULL,
    "START_DATE" date,
    "END_DATE" date,
    "INFO" character varying(255)
);

CREATE SEQUENCE "hari_libur_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "hari_libur_ID_seq" OWNED BY hari_libur."ID";

CREATE TABLE instansi (
    "ID" character varying(64) NOT NULL,
    "NAMA" character varying(255)
);

CREATE TABLE istri (
    "ID" bigint NOT NULL,
    "PNS" smallint,
    "NAMA" character varying(255),
    "TANGGAL_MENIKAH" date,
    "AKTE_NIKAH" character varying(255),
    "TANGGAL_MENINGGAL" date,
    "AKTE_MENINGGAL" character varying(255),
    "TANGGAL_CERAI" date,
    "AKTE_CERAI" character varying(255),
    "KARSUS" character varying(255),
    "STATUS" smallint,
    "HUBUNGAN" smallint,
    "PNS_ID" character varying(32),
    "NIP" character varying(32)
);

COMMENT ON COLUMN istri."PNS" IS '1=pns';

COMMENT ON COLUMN istri."STATUS" IS '1=menikah
2=cerai';

COMMENT ON COLUMN istri."HUBUNGAN" IS '1=istri, 2= suami';

CREATE SEQUENCE "istri_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "istri_ID_seq" OWNED BY istri."ID";

CREATE TABLE izin (
    "ID" integer NOT NULL,
    "NIP_PNS" character varying(18) NOT NULL,
    "NAMA" character varying(100),
    "JABATAN" character varying(255),
    "UNIT_KERJA" character varying(255),
    "MASA_KERJA_TAHUN" integer,
    "MASA_KERJA_BULAN" integer,
    "GAJI_POKOK" character varying(10),
    "KODE_IZIN" smallint NOT NULL,
    "DARI_TANGGAL" date,
    "SAMPAI_TANGGAL" date,
    "TAHUN" character varying(4),
    "JUMLAH" integer,
    "SATUAN" character varying(10),
    "KETERANGAN" character varying(255),
    "ALAMAT_SELAMA_CUTI" character varying(255),
    "TLP_SELAMA_CUTI" character varying(20),
    "TGL_DIBUAT" date,
    "LAMPIRAN_FILE" character varying(50),
    "SISA_CUTI_TAHUN_N2" integer,
    "SISA_CUTI_TAHUN_N1" integer,
    "SISA_CUTI_TAHUN_N" integer,
    "ANAK_KE" character varying(1),
    "NIP_ATASAN" character varying(25),
    "STATUS_ATASAN" integer,
    "CATATAN_ATASAN" character varying(255),
    "NIP_PYBMC" character varying(25),
    "STATUS_PYBMC" integer,
    "CATATAN_PYBMC" character varying(255),
    "NAMA_ATASAN" character varying(100),
    "NAMA_PYBMC" character varying(100),
    "TGL_PERKIRAAN_LAHIR" date,
    "TGL_ATASAN" date,
    "TGL_PPK" date,
    "NAMA_UNIT_KERJA" character varying(150),
    "ALASAN_CUTI" character varying(255),
    "SELAMA_JAM" character varying(20),
    "SELAMA_MENIT" character varying(20),
    "STATUS_PENGAJUAN" smallint DEFAULT 1,
    "NO_SURAT" character varying(100),
    "TUJUAN_JAUH" smallint,
    "TAMBAHAN_HARI" smallint,
    "LUAR_NEGERI" smallint DEFAULT 0,
    "TEXT_BASE64_SIGN" text,
    "IS_SIGNED" smallint,
    "DRAFT_BASE64_SIGN" text,
    created_at date,
    updated_at date,
    "JAM" time without time zone,
    ref uuid,
    source smallint DEFAULT 1,
    status_kirim smallint
);

COMMENT ON COLUMN izin."STATUS_PENGAJUAN" IS '1=menunggu persetujuan';

COMMENT ON COLUMN izin."LUAR_NEGERI" IS '0 = dalam negeri/null, 1= luar negeri';

COMMENT ON COLUMN izin."IS_SIGNED" IS '0';

COMMENT ON COLUMN izin.source IS '1=web,2=mobile';

COMMENT ON COLUMN izin.status_kirim IS '1=sudah pernah dikirim ke ekejadiran';

CREATE SEQUENCE "izin_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "izin_ID_seq" OWNED BY izin."ID";

CREATE TABLE izin_alasan (
    "ID" smallint NOT NULL,
    "ALASAN" character varying(255),
    "JENIS_CUTI" smallint
);

CREATE SEQUENCE "izin_alasan_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "izin_alasan_ID_seq" OWNED BY izin_alasan."ID";

CREATE TABLE izin_verifikasi (
    "ID" integer NOT NULL,
    "ID_PENGAJUAN" integer,
    "NIP_ATASAN" character varying(30),
    "STATUS_VERIFIKASI" smallint,
    "TANGGAL_VERIFIKASI" timestamp(6) without time zone,
    "ALASAN_DITOLAK" text
);

COMMENT ON COLUMN izin_verifikasi."STATUS_VERIFIKASI" IS '''id''=>1,''value''=>''Menunggu Persetujuan''),
			array(''id''=>2,''value''=>''Proses''),
			array(''id''=>3,''value''=>''Disetujui''),
			array(''id''=>4,''value''=>''Perubahan''),
			array(''id''=>5,''value''=>''Ditangguhkan''),
			array(''id''=>6,''value''=>''Ditolak''';

CREATE SEQUENCE "izin_verifikasi_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "izin_verifikasi_ID_seq" OWNED BY izin_verifikasi."ID";

CREATE SEQUENCE jabatan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE jabatan_id_seq OWNED BY jabatan.id;

CREATE TABLE jenis_arsip (
    "ID" integer NOT NULL,
    "NAMA_JENIS" character varying(255),
    "KETERANGAN" text,
    "KATEGORI_ARSIP" smallint
);

COMMENT ON COLUMN jenis_arsip."NAMA_JENIS" IS 'exa : ijazah SD, SK CPNS, SK PNS';

CREATE SEQUENCE jenis_arsip_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE jenis_arsip_id_seq OWNED BY jenis_arsip."ID";

CREATE TABLE jenis_diklat (
    id integer NOT NULL,
    bkn_id integer,
    jenis_diklat character varying(50),
    kode character varying(2),
    status smallint DEFAULT 1
);

CREATE SEQUENCE jenis_diklat_fungsional_id_seq
    START WITH 217
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE jenis_diklat_fungsional (
    "ID" bigint DEFAULT nextval('jenis_diklat_fungsional_id_seq'::regclass) NOT NULL,
    "NAMA" character varying(255) DEFAULT nextval('jenis_diklat_fungsional_id_seq'::regclass)
);

CREATE SEQUENCE jenis_diklat_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE jenis_diklat_id_seq OWNED BY jenis_diklat.id;

CREATE TABLE jenis_diklat_siasn (
    id bigint NOT NULL,
    jenis_diklat character varying
);

CREATE TABLE jenis_diklat_struktural (
    "ID" integer NOT NULL,
    "NAMA" character varying(255)
);

CREATE TABLE jenis_hukuman (
    "ID" character(2),
    "NAMA" character(100),
    "TINGKAT_HUKUMAN" character(1),
    "NAMA_TINGKAT_HUKUMAN" character(10)
);

CREATE TABLE jenis_izin (
    "ID" integer NOT NULL,
    "KODE" character varying(7) NOT NULL,
    "NAMA_IZIN" character varying(50) NOT NULL,
    "KETERANGAN" character varying(255),
    "PERSETUJUAN" character varying(100),
    "URUTAN" smallint,
    "STATUS" smallint DEFAULT 1,
    mobile smallint DEFAULT 0
);

CREATE SEQUENCE "jenis_izin_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "jenis_izin_ID_seq" OWNED BY jenis_izin."ID";

CREATE SEQUENCE rwt_jenis_jabatan_id_seq
    START WITH 3
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE jenis_jabatan (
    "ID" character varying(1) DEFAULT nextval('rwt_jenis_jabatan_id_seq'::regclass) NOT NULL,
    "NAMA" character varying(255)
);

CREATE TABLE jenis_kawin (
    "ID" character varying(255) NOT NULL,
    "NAMA" character varying(255)
);

CREATE TABLE jenis_kp (
    "ID" character varying(64) NOT NULL,
    "NAMA" character varying(255)
);

CREATE TABLE jenis_kursus (
    id integer NOT NULL,
    bkn_id character varying(255),
    kode_cepat character varying(10),
    jenis character varying(50),
    nama character varying(255)
);

CREATE SEQUENCE jenis_kursus_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE jenis_kursus_id_seq OWNED BY jenis_kursus.id;

CREATE TABLE jenis_pegawai (
    "ID" character varying(5) NOT NULL,
    "NAMA" character varying(200)
);

CREATE TABLE jenis_penghargaan (
    "ID" character(3) NOT NULL,
    "NAMA" character(100)
);

CREATE TABLE jenis_rumpun_diklat_siasn (
    id character varying NOT NULL,
    nama character varying,
    urusan character varying,
    pelayanan_dasar boolean,
    peraturan_id character varying,
    keterangan character varying
);

CREATE TABLE mst_jenis_satker (
    id_jenis smallint NOT NULL,
    nama_jenis_satker character varying(50)
);

CREATE SEQUENCE jenis_satker_id_jenis_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE jenis_satker_id_jenis_seq OWNED BY mst_jenis_satker.id_jenis;

CREATE TABLE kandidat_baperjakat (
    "ID" smallint NOT NULL,
    "NIP" character varying(32),
    "URUTAN_DEFAULT" smallint,
    "URUTAN_UPDATE" smallint,
    "TAHUN" character varying(4),
    "STATUS" smallint,
    "UNOR_ID" character varying(32),
    "STATUS_TAMBAHAN" smallint,
    "NILAI_ASSESMENT" real,
    "PANGGOL" character varying(255),
    "PENDIDIKAN" character varying(255),
    "HUKDIS" character varying(255),
    "UPDATE_BY" smallint,
    "UPDATE_DATE" date,
    "ID_PERIODE" integer,
    "JABATAN_ID" character varying(4),
    "NAMA_JABATAN" character varying(255),
    "TGL_PELANTIKAN" date,
    "NO_SK_PELANTIKAN" character varying(50),
    "KATEGORI" smallint,
    "STATUS_MENTERI" smallint
);

COMMENT ON COLUMN kandidat_baperjakat."STATUS" IS '1=diterima,0 tidak diterima';

COMMENT ON COLUMN kandidat_baperjakat."STATUS_TAMBAHAN" IS '1=admin, 2= sistem';

COMMENT ON COLUMN kandidat_baperjakat."KATEGORI" IS '1=rotasi,2=promosi';

CREATE SEQUENCE "kandidat_baperjakat_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "kandidat_baperjakat_ID_seq" OWNED BY kandidat_baperjakat."ID";

CREATE TABLE kategori_ds (
    id smallint NOT NULL,
    kategori_ds character varying(100)
);

CREATE SEQUENCE kategori_ds_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE kategori_ds_id_seq OWNED BY kategori_ds.id;

CREATE TABLE kategori_jenis_arsip (
    "ID" smallint NOT NULL,
    "KATEGORI" character varying(255)
);

CREATE SEQUENCE "kategori_jenis_arsip_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "kategori_jenis_arsip_ID_seq" OWNED BY kategori_jenis_arsip."ID";

CREATE TABLE kpkn (
    "ID" character varying(255) NOT NULL,
    "NAMA" character varying(255)
);

CREATE SEQUENCE kuota_jabatan_id_seq
    START WITH 6104
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE kuota_jabatan (
    "ID" bigint DEFAULT nextval('kuota_jabatan_id_seq'::regclass) NOT NULL,
    "KODE_UNIT_KERJA" character varying(50),
    "ID_JABATAN" character varying(32),
    "JUMLAH_PEMANGKU_JABATAN" smallint,
    "KETERANGAN" character varying(255),
    "FORMASI" character varying(50),
    "SKALA_PRIORITAS" character varying(50),
    "ID_JABATAN_PENYETARAAN" character varying(50),
    "PERATURAN" character varying(50),
    "KD_INTERNAL" character varying(50),
    kepmen_peta_jabatan character varying(100),
    nomor_kepmen_peta_jabatan character varying(100),
    tentang_kepmen_peta_jabatan character varying(200),
    aktif smallint
);

CREATE TABLE layanan (
    id bigint NOT NULL,
    layanan_tipe_id bigint,
    name character varying(255),
    keterangan character varying(255),
    expired_date date,
    _created_at timestamp(6) without time zone DEFAULT now() NOT NULL,
    active boolean
);

CREATE TABLE layanan_tipe (
    id bigint NOT NULL,
    name character varying(255),
    description character varying(255),
    active boolean
);

CREATE SEQUENCE layanan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE layanan_id_seq OWNED BY layanan_tipe.id;

CREATE SEQUENCE layanan_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE layanan_id_seq1 OWNED BY layanan.id;

CREATE TABLE layanan_usulan (
    id bigint NOT NULL,
    layanan_id bigint,
    pegawai_id bigint,
    created_at timestamp(6) without time zone,
    created_by bigint,
    status bigint,
    nip character varying(255),
    nama character varying(255),
    jabatan character varying(255),
    golongan_ruang character varying(255),
    unit_kerja character varying(255),
    satuan_kerja character varying(255),
    "F1" character varying(255),
    "F2" character varying(255),
    "F3" character varying(255),
    "F4" character varying(255),
    "F5" character varying(255),
    "F6" character varying(255),
    "F7" character varying(255),
    "F8" character varying(255),
    "F9" character varying(255),
    "F10" character varying(255),
    "F11" character varying(255),
    "F12" character varying(255),
    "F13" character varying(255),
    "F14" character varying(255),
    "F15" character varying(255),
    no_surat_pengantar character varying(255),
    file_surat_pengantar character varying(255)
);

CREATE SEQUENCE layanan_usulan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE layanan_usulan_id_seq OWNED BY layanan_usulan.id;

CREATE TABLE line_approval_izin (
    "ID" integer NOT NULL,
    "PNS_NIP" character varying(30),
    "NIP_ATASAN" character varying(30),
    "JENIS" smallint,
    "KETERANGAN_TAMBAHAN" character varying(200),
    "NAMA_ATASAN" character varying(100),
    "SEBAGAI" smallint
);

COMMENT ON COLUMN line_approval_izin."JENIS" IS '1=ATASAN LANGSUNG, 2 = PPK';

COMMENT ON COLUMN line_approval_izin."SEBAGAI" IS '1,2,3,4';

CREATE SEQUENCE "line_approval_izin_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "line_approval_izin_ID_seq" OWNED BY line_approval_izin."ID";

CREATE TABLE log_ds (
    "ID" integer NOT NULL,
    "ID_FILE" character varying(32),
    "NIK" character varying(30),
    "KETERANGAN" character varying(255),
    "CREATED_DATE" timestamp(0) without time zone,
    "CREATED_BY" integer,
    "STATUS" smallint,
    "PROSES_CRON" smallint DEFAULT 0
);

COMMENT ON COLUMN log_ds."STATUS" IS '1:gagal, 2:berhasil';

COMMENT ON COLUMN log_ds."PROSES_CRON" IS '0 = belum, 1 = sudah';

CREATE SEQUENCE "log_ds_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "log_ds_ID_seq" OWNED BY log_ds."ID";

CREATE TABLE log_request (
    id integer NOT NULL,
    url text NOT NULL,
    method character varying(10),
    params text,
    response_code integer,
    response text,
    created_at timestamp(6) without time zone DEFAULT now()
);

CREATE SEQUENCE log_request_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE log_request_id_seq OWNED BY log_request.id;

CREATE TABLE log_transaksi (
    "ID" integer NOT NULL,
    "NIP" character varying(30),
    "NAMA_KOMPUTER" character varying(30),
    "USER" character varying(30),
    "TGL_MODIFIKASI" date,
    "JAM_MODIFIKASI" time(6) without time zone,
    "YANG_DIUBAH" character varying(255),
    "MODULE" character varying(30)
);

CREATE SEQUENCE "log_transaksi_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "log_transaksi_ID_seq" OWNED BY log_transaksi."ID";

CREATE SEQUENCE login_attempts_id_seq
    START WITH 9073
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE login_attempts (
    ip_address character(40) NOT NULL,
    login character(50) NOT NULL,
    "time" time(6) without time zone,
    id integer DEFAULT nextval('login_attempts_id_seq'::regclass) NOT NULL
);

CREATE TABLE lokasi (
    "ID" character varying(255) NOT NULL,
    "KANREG_ID" character varying(255),
    "LOKASI_ID" character varying(255),
    "NAMA" character varying(255),
    "JENIS" character varying(255),
    "JENIS_KABUPATEN" character varying(255),
    "JENIS_DESA" character varying(255),
    "IBUKOTA" character varying(255)
);

CREATE VIEW mapping_unor_induk AS
 SELECT "ID",
    "NAMA_UNOR",
        CASE
            WHEN (("ESELON_2" IS NULL) OR (btrim(("ESELON_2")::text) = ''::text)) THEN "ESELON_1"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'universitas_dikti'::text) THEN "ESELON_2"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Multimedia Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Televisi Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Radio Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            ELSE "ESELON_2"
        END AS unor_induk,
        CASE
            WHEN (("ESELON_2" IS NULL) OR (btrim(("ESELON_2")::text) = ''::text)) THEN "NAMA_UNOR_ESELON_1"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'universitas_dikti'::text) THEN "NAMA_UNOR_ESELON_2"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Multimedia Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Televisi Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Radio Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            ELSE "NAMA_UNOR_ESELON_2"
        END AS nama_unor_induk
   FROM vw_unit_list t;

CREATE TABLE mst_peraturan_otk (
    id_peraturan smallint NOT NULL,
    no_peraturan character varying(100)
);

CREATE TABLE mst_templates (
    id integer NOT NULL,
    name character varying(255),
    template_file character varying(255)
);

CREATE SEQUENCE mst_templates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE mst_templates_id_seq OWNED BY mst_templates.id;

CREATE VIEW pns_aktif AS
 SELECT "ID",
        CASE
            WHEN ((masa_kerja[1] + bulan_swasta) >= 12) THEN ((masa_kerja[0] + tahun_swasta) + 1)
            ELSE (masa_kerja[0] + tahun_swasta)
        END AS masa_kerja_th,
        CASE
            WHEN ((masa_kerja[1] + bulan_swasta) >= 12) THEN ((masa_kerja[1] + bulan_swasta) - 12)
            ELSE (masa_kerja[1] + bulan_swasta)
        END AS masa_kerja_bl
   FROM ( SELECT pegawai."ID",
            get_masa_kerja_arr(pegawai."TMT_CPNS", ('now'::text)::date) AS masa_kerja,
            pegawai."MK_TAHUN_SWASTA" AS tahun_swasta,
            pegawai."MK_BULAN_SWASTA" AS bulan_swasta
           FROM pegawai
          WHERE ((pegawai.status_pegawai = 1) AND ((pegawai.terminated_date IS NULL) OR ((pegawai.terminated_date IS NOT NULL) AND (pegawai.terminated_date > ('now'::text)::date))))) temp;

CREATE MATERIALIZED VIEW mv_duk AS
 SELECT vw."NAMA_UNOR",
    pegawai."JENIS_JABATAN_ID",
    pegawai."JABATAN_ID",
    jabatan."NAMA_JABATAN",
    pegawai."NIP_LAMA",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    vw."ESELON_ID" AS vw_eselon_id,
    pegawai."GOL_ID",
    (((golongan."NAMA_PANGKAT")::text || ' '::text) || (golongan."NAMA")::text) AS golongan_text,
    'jabatanku'::text AS jabatan_text,
    pegawai."PNS_ID",
    (((date_part('year'::text, (now())::date) - date_part('year'::text, pegawai."TGL_LAHIR")) * (12)::double precision) + (date_part('month'::text, (now())::date) - date_part('month'::text, pegawai."TGL_LAHIR"))) AS bulan_usia,
    '#'::text AS separator,
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."KEDUDUKAN_HUKUM_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."KARTU_PEGAWAI",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."PENDIDIKAN",
    pegawai."TAHUN_LULUS",
    pegawai."KPKN_ID",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."ID",
    pegawai."PHOTO",
    pegawai."TMT_PENSIUN",
    pegawai."BUP",
    vw."ESELON_1",
    vw."ESELON_2",
    vw."ESELON_3",
    vw."ESELON_4",
    vw."NAMA_UNOR_ESELON_4",
    vw."NAMA_UNOR_ESELON_3",
    vw."NAMA_UNOR_ESELON_2",
    vw."NAMA_UNOR_ESELON_1"
   FROM ((((pns_aktif pa
     LEFT JOIN pegawai pegawai ON ((pa."ID" = pegawai."ID")))
     LEFT JOIN golongan ON (((pegawai."GOL_ID")::text = (golongan."ID")::text)))
     LEFT JOIN vw_unit_list vw ON (((vw."ID")::text = (pegawai."UNOR_ID")::text)))
     LEFT JOIN jabatan jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
  WHERE ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('78'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text]))
  ORDER BY pegawai."JENIS_JABATAN_ID", vw."ESELON_ID", vw."ESELON_1", vw."ESELON_2", vw."ESELON_3", vw."ESELON_4", pegawai."JABATAN_ID", vw."NAMA_UNOR_FULL", pegawai."GOL_ID" DESC, pegawai."TMT_GOLONGAN", pegawai."TMT_JABATAN", pegawai."TMT_CPNS", pegawai."TGL_LAHIR"
  WITH NO DATA;

CREATE MATERIALIZED VIEW mv_jml_unor_induk AS
 SELECT pegawai."UNOR_INDUK_ID",
    uk."JENIS_SATKER",
    uk."NAMA_UNOR",
    u."NAMA_UNOR" AS nama_unor_atas,
    count(pegawai."ID") AS jumlah
   FROM (((pegawai pegawai
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
     LEFT JOIN unitkerja uk ON (((uk."ID")::text = (pegawai."UNOR_INDUK_ID")::text)))
     LEFT JOIN unitkerja u ON (((u."ID")::text = (uk."DIATASAN_ID")::text)))
  WHERE (((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('88'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text, ('100'::character varying)::text])) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)) AND (pa."ID" IS NOT NULL))
  GROUP BY pegawai."UNOR_INDUK_ID", uk."JENIS_SATKER", uk."NAMA_UNOR", u."NAMA_UNOR"
  ORDER BY u."NAMA_UNOR", uk."NAMA_UNOR"
  WITH NO DATA;

CREATE TABLE tbl_file_ds (
    id_file character varying(200) NOT NULL,
    waktu_buat timestamp(0) without time zone,
    kategori character varying(100),
    teks_base64 text,
    id_pegawai_ttd character varying(255),
    is_signed smallint,
    nip_sk character varying(50),
    nomor_sk character varying(50),
    tgl_sk date,
    tmt_sk date,
    lokasi_file text,
    is_corrected smallint,
    catatan text,
    id_pegawai_korektor character varying(100),
    asal_surat_sk character varying(100),
    is_returned smallint,
    nama_pemilik_sk character varying(200),
    jabatan_pemilik_sk text,
    teks_base64_sign text,
    unit_kerja_pemilik_sk text,
    id integer NOT NULL,
    nip_pemroses character varying(50),
    ds_ok smallint,
    arsip character varying(50),
    "PNS_NONPNS" character varying(20),
    tmt_sampai_dengan date,
    telah_kirim smallint,
    halaman_ttd smallint DEFAULT 1,
    show_qrcode smallint DEFAULT 0,
    letak_ttd smallint DEFAULT 0,
    kode_unit_kerja_internal character varying(200),
    kode_jabatan_internal character varying(200),
    kelompok_jabatan character varying(200),
    tgl_tandatangan timestamp(6) without time zone,
    email_kirim character varying(200),
    sent_to_siasin character varying(100) DEFAULT 'n'::character varying,
    blockchain_issuer_id character varying,
    blockchain_image_url character varying,
    blockchain_hash character varying
);

COMMENT ON COLUMN tbl_file_ds.tmt_sampai_dengan IS 'khusus untuk Surat Perintah PLT/PLH';

COMMENT ON COLUMN tbl_file_ds.telah_kirim IS 'Jika 1, tampilkan di dikbudHR';

COMMENT ON COLUMN tbl_file_ds.halaman_ttd IS 'halaman diletakan tandataangan digital';

COMMENT ON COLUMN tbl_file_ds.show_qrcode IS '0/null : tidak tampilkan (seperti semula), 1 : tampilkan qrdari bssn';

COMMENT ON COLUMN tbl_file_ds.letak_ttd IS '1:tengah bawah, 2 : kiri Bawah 0: kanan bawah';

COMMENT ON COLUMN tbl_file_ds.kode_unit_kerja_internal IS 'untuk menampung nama unit kerja internal via kode';

COMMENT ON COLUMN tbl_file_ds.kode_jabatan_internal IS 'untuk menampung nama jabatan dengan kode jabatan internal';

COMMENT ON COLUMN tbl_file_ds.kelompok_jabatan IS 'khusus untuk keperluan laporan rekap';

COMMENT ON COLUMN tbl_file_ds.tgl_tandatangan IS 'untuk mengetahui tgl tandatangan';

COMMENT ON COLUMN tbl_file_ds.email_kirim IS 'Untuk menentukan alamat alternatif pengiriman dokumen';

CREATE MATERIALIZED VIEW mv_kategori_ds AS
 SELECT DISTINCT kategori AS kategori_ds
   FROM tbl_file_ds
  ORDER BY kategori
  WITH NO DATA;

CREATE TABLE pendidikan (
    "ID" character varying(255) NOT NULL,
    "TINGKAT_PENDIDIKAN_ID" character varying(255),
    "NAMA" character varying(255),
    "CEPAT_KODE" character varying(255)
);

CREATE MATERIALIZED VIEW mv_nominatif_pegawai AS
 SELECT pegawai."ID",
    btrim((pegawai."PNS_ID")::text) AS "PNS_ID",
    btrim((pegawai."GELAR_DEPAN")::text) AS "GELAR_DEPAN",
    btrim((pegawai."GELAR_BELAKANG")::text) AS "GELAR_BELAKANG",
    btrim((pegawai."NAMA")::text) AS "NAMA",
    btrim((pegawai."NIP_LAMA")::text) AS "NIP_LAMA",
    btrim((pegawai."NIP_BARU")::text) AS "NIP_BARU",
    btrim((pegawai."JENIS_KELAMIN")::text) AS "JENIS_KELAMIN",
    btrim((pegawai."TEMPAT_LAHIR_ID")::text) AS "TEMPAT_LAHIR_ID",
    pegawai."TGL_SK_CPNS",
    pegawai."TGL_LAHIR",
    date((pegawai."TGL_LAHIR" + ('1 year'::interval * (jabatan."PENSIUN")::double precision))) AS estimasi_pensiun,
    date_part('year'::text, age((pegawai."TGL_LAHIR")::timestamp with time zone)) AS age,
    pegawai."TMT_PENSIUN",
    btrim((pegawai."PENDIDIKAN_ID")::text) AS "PENDIDIKAN_ID",
    btrim((pendidikan."NAMA")::text) AS "NAMA_PENDIDIKAN",
    btrim((agama."NAMA")::text) AS "NAMA_AGAMA",
    btrim((jenis_jabatan."NAMA")::text) AS "JENIS_JABATAN",
    jabatan."KELAS",
    btrim((jabatan."NAMA_JABATAN")::text) AS "NAMA_JABATAN",
    btrim((golongan."NAMA")::text) AS "NAMA_GOLONGAN",
    btrim((golongan."NAMA_PANGKAT")::text) AS "NAMA_PANGKAT",
    pegawai."GOL_ID",
    pegawai."TMT_GOLONGAN",
    btrim((vw."NAMA_UNOR_ESELON_4")::text) AS "NAMA_UNOR_ESELON_4",
    btrim((vw."NAMA_UNOR_ESELON_3")::text) AS "NAMA_UNOR_ESELON_3",
    btrim((vw."NAMA_UNOR_ESELON_2")::text) AS "NAMA_UNOR_ESELON_2",
    btrim((vw."NAMA_UNOR_ESELON_1")::text) AS "NAMA_UNOR_ESELON_1",
    btrim((vw."ID")::text) AS "ID_UNOR",
    btrim((vw."ESELON_1")::text) AS "ESELON_1",
    btrim((vw."ESELON_2")::text) AS "ESELON_2",
    btrim((vw."ESELON_3")::text) AS "ESELON_3",
    btrim((vw."ESELON_4")::text) AS "ESELON_4",
    btrim((vw."ESELON_ID")::text) AS "ESELON_ID",
    vw."JENIS_SATKER",
    pegawai."TK_PENDIDIKAN",
    jabatan."KATEGORI_JABATAN",
    btrim((kedudukan_hukum."NAMA")::text) AS "KEDUDUKAN_HUKUM_NAMA",
    btrim((pegawai."NOMOR_SK_CPNS")::text) AS "NOMOR_SK_CPNS",
    pegawai."TMT_CPNS",
    btrim((uk."NAMA_UNOR")::text) AS nama_satker,
    btrim((pegawai."NIK")::text) AS "NIK",
    btrim((pegawai."NOMOR_HP")::text) AS "NOMOR_HP",
    btrim((pegawai."NOMOR_DARURAT")::text) AS "NOMOR_DARURAT",
    btrim((pegawai."EMAIL")::text) AS "EMAIL",
    btrim((pegawai."EMAIL_DIKBUD")::text) AS "EMAIL_DIKBUD",
    vw."NAMA_UNOR_FULL",
    lokasi."NAMA" AS "TEMPAT_LAHIR_NAMA"
   FROM ((((((((((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
     LEFT JOIN unitkerja uk ON (((uk."ID")::text = (vw."UNOR_INDUK")::text)))
     LEFT JOIN golongan ON ((pegawai."GOL_ID" = golongan."ID")))
     LEFT JOIN lokasi ON (((lokasi."ID")::text = (pegawai."TEMPAT_LAHIR_ID")::text)))
     LEFT JOIN pendidikan ON (((pendidikan."ID")::text = (pegawai."PENDIDIKAN_ID")::text)))
     LEFT JOIN agama ON ((agama."ID" = pegawai."AGAMA_ID")))
     LEFT JOIN kedudukan_hukum ON (((kedudukan_hukum."ID")::text = (pegawai."KEDUDUKAN_HUKUM_ID")::text)))
     LEFT JOIN jabatan ON ((pegawai."JABATAN_INSTANSI_REAL_ID" = (jabatan."KODE_JABATAN")::bpchar)))
     LEFT JOIN jenis_jabatan ON (((jenis_jabatan."ID")::text = (jabatan."JENIS_JABATAN")::text)))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('78'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text])) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)))
  ORDER BY (btrim((pegawai."NAMA")::text))
  WITH NO DATA;

CREATE TABLE tkpendidikan (
    "ID" character varying(255) NOT NULL,
    "GOLONGAN_ID" character varying(255),
    "NAMA" character varying(255),
    "GOLONGAN_AWAL_ID" character varying(255),
    "DELETED" smallint,
    "ABBREVIATION" character varying(255),
    "TINGKAT" smallint
);

CREATE MATERIALIZED VIEW mv_pegawai AS
 SELECT btrim((pegawai."ID")::text) AS "ID",
    btrim((pegawai."PNS_ID")::text) AS "PNS_ID",
    btrim((pegawai."NIP_LAMA")::text) AS "NIP_LAMA",
    btrim((pegawai."NIP_BARU")::text) AS "NIP_BARU",
    btrim((pegawai."NAMA")::text) AS "NAMA",
    btrim((pegawai."GELAR_DEPAN")::text) AS "GELAR_DEPAN",
    btrim((pegawai."GELAR_BELAKANG")::text) AS "GELAR_BELAKANG",
    btrim((pegawai."TEMPAT_LAHIR_ID")::text) AS "TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    btrim((pegawai."JENIS_KELAMIN")::text) AS "JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    btrim((pegawai."JENIS_KAWIN_ID")::text) AS "JENIS_KAWIN_ID",
    btrim((pegawai."NIK")::text) AS "NIK",
    btrim((pegawai."NOMOR_DARURAT")::text) AS "NOMOR_DARURAT",
    btrim((pegawai."NOMOR_HP")::text) AS "NOMOR_HP",
    btrim((pegawai."EMAIL")::text) AS "EMAIL",
    btrim((pegawai."ALAMAT")::text) AS "ALAMAT",
    btrim((pegawai."NPWP")::text) AS "NPWP",
    btrim((pegawai."BPJS")::text) AS "BPJS",
    btrim((pegawai."JENIS_PEGAWAI_ID")::text) AS "JENIS_PEGAWAI_ID",
    btrim((pegawai."KEDUDUKAN_HUKUM_ID")::text) AS "KEDUDUKAN_HUKUM_ID",
    btrim((pegawai."STATUS_CPNS_PNS")::text) AS "STATUS_CPNS_PNS",
    btrim((pegawai."KARTU_PEGAWAI")::text) AS "KARTU_PEGAWAI",
    btrim((pegawai."NOMOR_SK_CPNS")::text) AS "NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    btrim((pegawai."GOL_AWAL_ID")::text) AS "GOL_AWAL_ID",
    pegawai."GOL_ID",
    pegawai."TMT_GOLONGAN",
    btrim((pegawai."MK_TAHUN")::text) AS "MK_TAHUN",
    btrim((pegawai."MK_BULAN")::text) AS "MK_BULAN",
    btrim((pegawai."JENIS_JABATAN_IDx")::text) AS "JENIS_JABATAN_IDx",
    btrim((pegawai."JABATAN_ID")::text) AS "JABATAN_ID",
    pegawai."TMT_JABATAN",
    btrim((pegawai."PENDIDIKAN_ID")::text) AS "PENDIDIKAN_ID",
    btrim((pendidikan."NAMA")::text) AS "NAMA_PENDIDIKAN",
    btrim((tkpendidikan."NAMA")::text) AS "TINGKAT_PENDIDIKAN_NAMA",
    btrim((pegawai."TAHUN_LULUS")::text) AS "TAHUN_LULUS",
    btrim((pegawai."KPKN_ID")::text) AS "KPKN_ID",
    btrim((pegawai."LOKASI_KERJA_ID")::text) AS "LOKASI_KERJA_ID",
    btrim((pegawai."UNOR_ID")::text) AS "UNOR_ID",
    btrim((pegawai."UNOR_INDUK_ID")::text) AS "UNOR_INDUK_ID",
    btrim((pegawai."INSTANSI_INDUK_ID")::text) AS "INSTANSI_INDUK_ID",
    btrim((pegawai."INSTANSI_KERJA_ID")::text) AS "INSTANSI_KERJA_ID",
    btrim((pegawai."SATUAN_KERJA_INDUK_ID")::text) AS "SATUAN_KERJA_INDUK_ID",
    btrim((pegawai."SATUAN_KERJA_KERJA_ID")::text) AS "SATUAN_KERJA_KERJA_ID",
    btrim((pegawai."GOLONGAN_DARAH")::text) AS "GOLONGAN_DARAH",
    btrim((pegawai."PHOTO")::text) AS "PHOTO",
    pegawai."TMT_PENSIUN",
    btrim((pegawai."LOKASI_KERJA")::text) AS "LOKASI_KERJA",
    btrim((pegawai."JML_ISTRI")::text) AS "JML_ISTRI",
    btrim((pegawai."JML_ANAK")::text) AS "JML_ANAK",
    btrim((pegawai."NO_SURAT_DOKTER")::text) AS "NO_SURAT_DOKTER",
    pegawai."TGL_SURAT_DOKTER",
    btrim((pegawai."NO_BEBAS_NARKOBA")::text) AS "NO_BEBAS_NARKOBA",
    pegawai."TGL_BEBAS_NARKOBA",
    btrim((pegawai."NO_CATATAN_POLISI")::text) AS "NO_CATATAN_POLISI",
    pegawai."TGL_CATATAN_POLISI",
    btrim((pegawai."AKTE_KELAHIRAN")::text) AS "AKTE_KELAHIRAN",
    btrim((pegawai."STATUS_HIDUP")::text) AS "STATUS_HIDUP",
    btrim((pegawai."AKTE_MENINGGAL")::text) AS "AKTE_MENINGGAL",
    pegawai."TGL_MENINGGAL",
    btrim((pegawai."NO_ASKES")::text) AS "NO_ASKES",
    btrim((pegawai."NO_TASPEN")::text) AS "NO_TASPEN",
    pegawai."TGL_NPWP",
    btrim((pegawai."TEMPAT_LAHIR")::text) AS "TEMPAT_LAHIR",
    btrim((pegawai."PENDIDIKAN")::text) AS "PENDIDIKAN",
    btrim((pegawai."TK_PENDIDIKAN")::text) AS "TK_PENDIDIKAN",
    btrim((pegawai."TEMPAT_LAHIR_NAMA")::text) AS "TEMPAT_LAHIR_NAMA",
    btrim((pegawai."JENIS_JABATAN_NAMA")::text) AS "JENIS_JABATAN_NAMA",
    btrim((pegawai."JABATAN_NAMA")::text) AS "JABATAN_NAMA",
    btrim((pegawai."KPKN_NAMA")::text) AS "KPKN_NAMA",
    btrim((pegawai."INSTANSI_INDUK_NAMA")::text) AS "INSTANSI_INDUK_NAMA",
    btrim((pegawai."INSTANSI_KERJA_NAMA")::text) AS "INSTANSI_KERJA_NAMA",
    btrim((pegawai."SATUAN_KERJA_INDUK_NAMA")::text) AS "SATUAN_KERJA_INDUK_NAMA",
    btrim((pegawai."SATUAN_KERJA_NAMA")::text) AS "SATUAN_KERJA_NAMA",
    btrim((pegawai."JABATAN_INSTANSI_ID")::text) AS "JABATAN_INSTANSI_ID",
    btrim((pegawai."JABATAN_INSTANSI_NAMA")::text) AS "JABATAN_INSTANSI_NAMA",
    pegawai."JENIS_JABATAN_ID",
    pegawai.terminated_date,
    pegawai.status_pegawai,
    btrim((pegawai."JABATAN_PPNPN")::text) AS "JABATAN_PPNPN",
    btrim((jr."NAMA_JABATAN")::text) AS "NAMA_JABATAN_REAL",
    btrim((jr."KATEGORI_JABATAN")::text) AS "KATEGORI_JABATAN_REAL",
    jr."JENIS_JABATAN" AS "JENIS_JABATAN_REAL",
    pegawai."CREATED_DATE",
    pegawai."CREATED_BY",
    pegawai."UPDATED_DATE",
    pegawai."UPDATED_BY",
    btrim((pegawai."EMAIL_DIKBUD")::text) AS "EMAIL_DIKBUD",
    btrim((pegawai."KODECEPAT")::text) AS "KODECEPAT",
    vw."NAMA_UNOR_FULL",
    btrim((golongan."NAMA")::text) AS "NAMA_GOLONGAN",
    btrim((golongan."NAMA_PANGKAT")::text) AS "NAMA_PANGKAT",
    btrim((vw."NAMA_UNOR_ESELON_4")::text) AS "NAMA_UNOR_ESELON_4",
    btrim((vw."NAMA_UNOR_ESELON_3")::text) AS "NAMA_UNOR_ESELON_3",
    btrim((vw."NAMA_UNOR_ESELON_2")::text) AS "NAMA_UNOR_ESELON_2",
    btrim((vw."NAMA_UNOR_ESELON_1")::text) AS "NAMA_UNOR_ESELON_1",
    btrim((un."NAMA_UNOR")::text) AS "UNOR_INDUK_NAMA",
    btrim((vw."ESELON_1")::text) AS "ESELON_1",
    btrim((vw."ESELON_2")::text) AS "ESELON_2",
    btrim((vw."ESELON_3")::text) AS "ESELON_3",
    btrim((vw."ESELON_4")::text) AS "ESELON_4",
    btrim((vw."ESELON_ID")::text) AS "ESELON_ID",
    btrim((kedudukan_hukum."NAMA")::text) AS "KEDUDUKAN_HUKUM_NAMA",
    pa."ID" AS "PNS_AKTIF_ID",
    jr."KELAS" AS "KELAS_JABATAN",
    btrim((agama."NAMA")::text) AS "NAMA_AGAMA"
   FROM ((((((((((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN golongan ON ((pegawai."GOL_ID" = golongan."ID")))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
     LEFT JOIN jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
     LEFT JOIN jabatan jr ON ((pegawai."JABATAN_INSTANSI_REAL_ID" = (jr."KODE_JABATAN")::bpchar)))
     LEFT JOIN pendidikan ON (((pegawai."PENDIDIKAN_ID")::bpchar = (pendidikan."ID")::bpchar)))
     LEFT JOIN tkpendidikan ON (((tkpendidikan."ID")::bpchar = (pendidikan."TINGKAT_PENDIDIKAN_ID")::bpchar)))
     LEFT JOIN unitkerja un ON (((vw."UNOR_INDUK")::text = (un."ID")::text)))
     LEFT JOIN agama ON ((agama."ID" = pegawai."AGAMA_ID")))
     LEFT JOIN kedudukan_hukum ON (((kedudukan_hukum."ID")::text = (pegawai."KEDUDUKAN_HUKUM_ID")::text)))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '14'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '52'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '66'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '67'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '77'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '78'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '98'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '99'::text) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)))
  ORDER BY pegawai."NAMA"
  WITH NO DATA;

CREATE TABLE pegawai_atasan (
    "ID" integer NOT NULL,
    "PNS_NIP" character varying(18),
    "NIP_ATASAN" character varying(18),
    "PPK" character varying(18),
    "KETERANGAN_TAMBAHAN" character varying(255),
    "NAMA_ATASAN" character varying(100),
    "NAMA_PPK" character varying(100)
);

CREATE MATERIALIZED VIEW mv_pegawai_cuti AS
 SELECT pegawai."ID",
    pegawai."NIP_BARU",
    pegawai."PNS_ID",
    btrim((pegawai."GELAR_DEPAN")::text) AS "GELAR_DEPAN",
    btrim((pegawai."NAMA")::text) AS "NAMA",
    btrim((pegawai."GELAR_BELAKANG")::text) AS "GELAR_BELAKANG",
    btrim((pegawai."UNOR_INDUK_ID")::text) AS "UNOR_INDUK_ID",
    btrim((pegawai."UNOR_ID")::text) AS "UNOR_ID",
    golongan."NAMA" AS "NAMA_GOLONGAN",
    golongan."NAMA_PANGKAT",
    jabatan."KODE_JABATAN",
    jabatan."NAMA_JABATAN",
    jabatan."KATEGORI_JABATAN",
    btrim((pt."NIP_ATASAN")::text) AS "NIP_ATASAN",
    btrim((pt."NAMA_ATASAN")::text) AS "NAMA_ATASAN",
    btrim((pt."PPK")::text) AS "PPK",
    btrim((pt."NAMA_PPK")::text) AS "NAMA_PPK",
    pt."ID" AS "ID_PEGAWAI_ATASAN",
    pt."KETERANGAN_TAMBAHAN",
    vw."UNOR_INDUK",
    vw."NAMA_UNOR_FULL"
   FROM (((((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN golongan ON ((pegawai."GOL_ID" = golongan."ID")))
     LEFT JOIN jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
     LEFT JOIN pegawai_atasan pt ON (((pegawai."NIP_BARU")::text = (pt."PNS_NIP")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '99'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '66'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '52'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '20'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '04'::text) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)))
  ORDER BY vw."UNOR_INDUK"
  WITH NO DATA;

CREATE MATERIALIZED VIEW mv_pegawai_layanan AS
 SELECT "ID" AS id,
    "NIP_BARU" AS nip,
    "NAMA" AS nama,
    "UNOR_ID" AS unitid,
    "GOL_ID" AS golid,
    "JABATAN_INSTANSI_ID" AS jabid,
    "KEDUDUKAN_HUKUM_ID" AS status_hukum,
    "TMT_CPNS" AS tmtcpns,
    "TMT_PNS" AS tmtpns,
    "MK_TAHUN" AS thn,
    "MK_BULAN" AS bln,
    "TMT_PENSIUN" AS tmtpensiun,
    "GELAR_DEPAN" AS gelar_depan,
    "GELAR_BELAKANG" AS gelar_belakang,
    "JENIS_KELAMIN" AS jeniskelamin,
    "TEMPAT_LAHIR" AS tempatlahir,
    "TGL_LAHIR" AS tgllahir
   FROM pegawai pegawai
  WITH NO DATA;

CREATE TABLE rwt_assesmen (
    "ID" integer NOT NULL,
    "PNS_ID" character(255),
    "PNS_NIP" character(25),
    "TAHUN" character(4),
    "FILE_UPLOAD" character(100),
    "NILAI" real,
    "NILAI_KINERJA" real,
    "TAHUN_PENILAIAN_ID" character varying(32),
    "TAHUN_PENILAIAN_TITLE" character varying(20) NOT NULL,
    "FULLNAME" character varying(255),
    "POSISI_ID" character varying(20),
    "UNIT_ORG_ID" character varying(50),
    "NAMA_UNOR" character varying(200),
    "SARANPENGEMBANGAN" text,
    "FILE_UPLOAD_FB_POTENSI" character varying(255),
    "FILE_UPLOAD_LENGKAP_PT" character varying(255),
    "FILE_UPLOAD_FB_PT" character varying(255),
    "FILE_UPLOAD_EXISTS" smallint DEFAULT 0,
    "SATKER_ID" character varying
);

COMMENT ON COLUMN rwt_assesmen."TAHUN" IS 'tahun_penilaian_awal';

COMMENT ON COLUMN rwt_assesmen."FILE_UPLOAD" IS 'file laporan lengkap potensi';

COMMENT ON COLUMN rwt_assesmen."NILAI" IS 'Nilai Potensi';

CREATE MATERIALIZED VIEW mv_riwayat_asesmen AS
 SELECT btrim((rwt_assesmen."PNS_NIP")::text) AS "PNS_NIP",
    btrim((rwt_assesmen."TAHUN")::text) AS "TAHUN",
    rwt_assesmen."NILAI",
    btrim((rwt_assesmen."FILE_UPLOAD")::text) AS "FILE_UPLOAD_FB_POTENSI",
    btrim((rwt_assesmen."FILE_UPLOAD")::text) AS "FILE_UPLOAD_LENGKAP_PT",
    btrim((rwt_assesmen."FILE_UPLOAD")::text) AS "FILE_UPLOAD_FB_PT"
   FROM rwt_assesmen
  WHERE ((btrim((rwt_assesmen."TAHUN")::text) = ('2019'::bpchar)::text) AND (rwt_assesmen."FILE_UPLOAD" ~~* '%a_p%'::text) AND (rwt_assesmen."FILE_UPLOAD_EXISTS" = '1'::smallint))
UNION ALL
 SELECT btrim((rwt_assesmen."PNS_NIP")::text) AS "PNS_NIP",
    rwt_assesmen."TAHUN",
    rwt_assesmen."NILAI",
    btrim((rwt_assesmen."FILE_UPLOAD_FB_POTENSI")::text) AS "FILE_UPLOAD_FB_POTENSI",
    btrim((rwt_assesmen."FILE_UPLOAD_LENGKAP_PT")::text) AS "FILE_UPLOAD_LENGKAP_PT",
    btrim((rwt_assesmen."FILE_UPLOAD_FB_PT")::text) AS "FILE_UPLOAD_FB_PT"
   FROM rwt_assesmen
  WHERE (btrim((rwt_assesmen."TAHUN")::text) <> ('2019'::bpchar)::text)
  WITH NO DATA;

CREATE MATERIALIZED VIEW mv_unit_list_all AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(32) AS "ESELON_1",
    (x.eselon[2])::character varying(32) AS "ESELON_2",
    (x.eselon[3])::character varying(32) AS "ESELON_3",
    (x.eselon[4])::character varying(32) AS "ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WITH NO DATA;

CREATE VIEW nama_unit AS
 SELECT "ID",
    "NAMA_UNOR",
    ( SELECT unitkerja."NAMA_UNOR"
           FROM unitkerja unitkerja
          WHERE ((unitkerja."ID")::text = (a."UNOR_INDUK")::text)) AS satker,
    ( SELECT unitkerja."NAMA_UNOR"
           FROM unitkerja unitkerja
          WHERE ((unitkerja."ID")::text = (a."ESELON_4")::text)) AS es4,
    ( SELECT unitkerja."NAMA_UNOR"
           FROM unitkerja unitkerja
          WHERE ((unitkerja."ID")::text = (a."ESELON_3")::text)) AS es3,
    ( SELECT unitkerja."NAMA_UNOR"
           FROM unitkerja unitkerja
          WHERE ((unitkerja."ID")::text = (a."ESELON_2")::text)) AS es2,
    ( SELECT unitkerja."NAMA_UNOR"
           FROM unitkerja unitkerja
          WHERE ((unitkerja."ID")::text = (a."ESELON_1")::text)) AS es1
   FROM unitkerja a;

CREATE TABLE nip_pejabat (
    "NIP" character varying(18) NOT NULL,
    id bigint NOT NULL
);

CREATE SEQUENCE nip_pejabat_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE nip_pejabat_id_seq OWNED BY nip_pejabat.id;

CREATE TABLE orang_tua (
    "ID" integer NOT NULL,
    "HUBUNGAN" smallint,
    "ALAMAT" text,
    "NO_TLP" character varying(30),
    "NO_HP" character varying(50),
    "STATUS_PERKAWINAN" character varying(20),
    "AKTE_KELAHIRAN" character varying(255),
    "STATUS_HIDUP" smallint,
    "AKTE_MENINGGAL" character varying(255),
    "TGL_MENINGGAL" date,
    "NO_NPWP" character varying(255),
    "TANGGAL_NPWP" date,
    "NAMA" character varying(255),
    "GELAR_DEPAN" character varying(20),
    "GELAR_BELAKANG" character varying(50),
    "TEMPAT_LAHIR" character varying(255),
    "TANGGAL_LAHIR" character varying(255),
    "JENIS_KELAMIN" character varying(20),
    "AGAMA" character varying(2),
    "EMAIL" character varying(255),
    "JENIS_DOKUMEN_ID" character varying(10),
    "NO_DOKUMEN_ID" character varying(50),
    "FOTO" character varying(255),
    "KODE" smallint,
    "NIP" character varying(32),
    "PNS_ID" character varying(32)
);

COMMENT ON COLUMN orang_tua."KODE" IS '1=AYAH
2=IBU';

CREATE SEQUENCE "orang_tua_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "orang_tua_ID_seq" OWNED BY orang_tua."ID";

CREATE VIEW organisasi_ropeg AS
 WITH RECURSIVE org AS (
         SELECT unitkerja."ID" AS idorg,
            unitkerja."DIATASAN_ID" AS idorgparent,
            unitkerja."NAMA_UNOR" AS orgname,
            unitkerja."PEMIMPIN_PNS_ID" AS pemimpinpnsid
           FROM unitkerja
          WHERE ((unitkerja."ID")::text = '8ae483a8641f817901641fce97d21d1b'::text)
        UNION
         SELECT e."ID",
            e."DIATASAN_ID",
            e."NAMA_UNOR",
            e."PEMIMPIN_PNS_ID"
           FROM (unitkerja e
             JOIN org s ON (((s.idorg)::text = (e."DIATASAN_ID")::text)))
        )
 SELECT idorg,
    idorgparent,
    orgname,
    pemimpinpnsid
   FROM org x
  ORDER BY idorg;

CREATE SEQUENCE "pegawai_atasan_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "pegawai_atasan_ID_seq" OWNED BY pegawai_atasan."ID";

CREATE SEQUENCE pegawai_bkn_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE pegawai_bkn (
    "ID" integer DEFAULT nextval('pegawai_bkn_id'::regclass) NOT NULL,
    "PNS_ID" character varying(32) NOT NULL,
    "NIP_LAMA" character varying(9),
    "NIP_BARU" character varying(18),
    "NAMA" character varying(255),
    "GELAR_DEPAN" character varying(60),
    "GELAR_BELAKANG" character varying(60),
    "TEMPAT_LAHIR_ID" character varying(50),
    "TGL_LAHIR" date,
    "JENIS_KELAMIN" character varying(10),
    "AGAMA_ID" integer,
    "JENIS_KAWIN_ID" character varying(255),
    "NIK" character varying(255),
    "NOMOR_DARURAT" character varying(255),
    "NOMOR_HP" character varying(60),
    "EMAIL" character varying(255),
    "ALAMAT" character varying(255),
    "NPWP" character varying(255),
    "BPJS" character varying(50),
    "JENIS_PEGAWAI_ID" character varying(50),
    "KEDUDUKAN_HUKUM_ID" character varying(32),
    "STATUS_CPNS_PNS" character varying(20),
    "KARTU_PEGAWAI" character varying(30),
    "NOMOR_SK_CPNS" character varying(60),
    "TGL_SK_CPNS" date,
    "TMT_CPNS" date,
    "TMT_PNS" date,
    "GOL_AWAL_ID" character varying(32),
    "GOL_ID" integer,
    "TMT_GOLONGAN" date,
    "MK_TAHUN" character varying(20),
    "MK_BULAN" character varying(20),
    "JENIS_JABATAN_IDx" character varying(32),
    "JABATAN_ID" character varying(32),
    "TMT_JABATAN" date,
    "PENDIDIKAN_ID" character varying(32),
    "TAHUN_LULUS" character varying(20),
    "KPKN_ID" character varying(32),
    "LOKASI_KERJA_ID" character varying(32),
    "UNOR_ID" character varying(32),
    "UNOR_INDUK_ID" character varying(32),
    "INSTANSI_INDUK_ID" character varying(32),
    "INSTANSI_KERJA_ID" character varying(32),
    "SATUAN_KERJA_INDUK_ID" character varying(32),
    "SATUAN_KERJA_KERJA_ID" character varying(32),
    "GOLONGAN_DARAH" character varying(20),
    "PHOTO" character varying(100),
    "TMT_PENSIUN" date,
    "LOKASI_KERJA" character(200),
    "JML_ISTRI" character(1),
    "JML_ANAK" character(1),
    "NO_SURAT_DOKTER" character(100),
    "TGL_SURAT_DOKTER" date,
    "NO_BEBAS_NARKOBA" character(100),
    "TGL_BEBAS_NARKOBA" date,
    "NO_CATATAN_POLISI" character(100),
    "TGL_CATATAN_POLISI" date,
    "AKTE_KELAHIRAN" character(50),
    "STATUS_HIDUP" character(15),
    "AKTE_MENINGGAL" character(50),
    "TGL_MENINGGAL" date,
    "NO_ASKES" character(50),
    "NO_TASPEN" character(50),
    "TGL_NPWP" date,
    "TEMPAT_LAHIR" character(200),
    "PENDIDIKAN" character(165),
    "TK_PENDIDIKAN" character(3),
    "TEMPAT_LAHIR_NAMA" character(200),
    "JENIS_JABATAN_NAMA" character(200),
    "JABATAN_NAMA" character(254),
    "KPKN_NAMA" character(255),
    "INSTANSI_INDUK_NAMA" character(100),
    "INSTANSI_KERJA_NAMA" character(160),
    "SATUAN_KERJA_INDUK_NAMA" character(170),
    "SATUAN_KERJA_NAMA" character(155),
    "JABATAN_INSTANSI_ID" character(15),
    "BUP" smallint DEFAULT 58,
    "JABATAN_INSTANSI_NAMA" character varying(512) DEFAULT NULL::character varying,
    "JENIS_JABATAN_ID" integer
);

CREATE TABLE pengajuan_tubel (
    "ID" integer NOT NULL,
    "NIP" character varying(30),
    "NOMOR_USUL" character varying(20),
    "TANGGAL_USUL" date,
    "UNIVERSITAS" character varying(100),
    "FAKULTAS" character varying(100),
    "PRODI" character varying(100),
    "BEASISWA" integer,
    "PEMBERI_BEASISWA" character varying(100),
    "JENJANG" character varying(5),
    "NEGARA" character varying(50),
    "STATUS" integer DEFAULT 1,
    "ALASAN_DITOLAK" character varying(255),
    "MULAI_BELAJAR" date,
    "AKHIR_BELAJAR" date
);

CREATE SEQUENCE "pengajuan_tubel_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "pengajuan_tubel_ID_seq" OWNED BY pengajuan_tubel."ID";

CREATE SEQUENCE peraturan_otk_id_peraturan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE peraturan_otk_id_peraturan_seq OWNED BY mst_peraturan_otk.id_peraturan;

CREATE TABLE perkiraan_kpo (
    id bigint NOT NULL,
    nip character varying(255),
    status smallint,
    alasan text,
    layanan_id bigint,
    nama character varying(255),
    birth_place character varying(255),
    birth_date date,
    last_edu character varying(255),
    o_gol_ruang character varying(255),
    o_gol_tmt date,
    o_masakerja_thn smallint,
    o_masakerja_bln smallint,
    o_gapok double precision,
    o_jabatan character varying(255),
    o_tmt_jabatan date,
    n_gol_ruang character varying(255),
    n_gol_tmt date,
    n_masakerja_thn smallint,
    n_masakerja_bln smallint,
    n_gapok double precision,
    n_jabatan character varying(255),
    n_tmt_jabatan date,
    unit_kerja character varying(255),
    unit_kerja_induk character varying(255),
    kantor_pembayaran character varying(255),
    tahun_lulus smallint,
    no_surat_pengantar character varying(255),
    no_surat_pengantar_es1 character varying(255)
);

COMMENT ON COLUMN perkiraan_kpo.status IS 'tms/ms';

CREATE SEQUENCE perkiraan_kpo_documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE perkiraan_kpo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE perkiraan_kpo_id_seq OWNED BY perkiraan_kpo.id;

CREATE SEQUENCE perkiraan_ppo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE perkiraan_ppo (
    id bigint DEFAULT nextval('perkiraan_ppo_id_seq'::regclass) NOT NULL,
    nip character varying(255),
    status smallint DEFAULT 1,
    alasan text,
    layanan_id bigint,
    nama character varying(255),
    birth_place character varying(255),
    birth_date date,
    last_edu character varying(255),
    o_gol_ruang character varying(255),
    o_gol_tmt date,
    o_masakerja_thn smallint,
    o_masakerja_bln smallint,
    o_gapok double precision,
    o_jabatan character varying(255),
    o_tmt_jabatan date,
    n_gol_ruang character varying(255),
    n_gol_tmt date,
    n_masakerja_thn smallint,
    n_masakerja_bln smallint,
    n_gapok double precision,
    n_jabatan character varying(255),
    n_tmt_jabatan date,
    unit_kerja character varying(255),
    unit_kerja_induk character varying(255),
    kantor_pembayaran character varying(255),
    tahun_lulus smallint,
    no_surat_pengantar character varying(255),
    bup smallint,
    n_jabatan_id character varying
);

COMMENT ON COLUMN perkiraan_ppo.status IS 'tms/ms';

CREATE TABLE perkiraan_usulan_log (
    id bigint NOT NULL,
    usulan_id bigint,
    _created_at timestamp(6) without time zone DEFAULT now(),
    _created_by character varying(255),
    status character varying(255),
    alasan character varying(255)
);

CREATE SEQUENCE perkiraan_usulan_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE perkiraan_usulan_log_id_seq OWNED BY perkiraan_usulan_log.id;

CREATE SEQUENCE permissions_permission_id_seq
    START WITH 219
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE permissions (
    permission_id bigint DEFAULT nextval('permissions_permission_id_seq'::regclass) NOT NULL,
    name character varying(100),
    description character varying(255),
    status character varying(20)
);

CREATE TABLE peta_jabatan_permen (
    id smallint NOT NULL,
    permen character varying(50)
);

CREATE SEQUENCE peta_jabatan_permen_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE peta_jabatan_permen_id_seq OWNED BY peta_jabatan_permen.id;

CREATE TABLE pindah_unit (
    "ID" integer NOT NULL,
    "NIP" character varying(32) NOT NULL,
    "SURAT_PERMOHONAN_PINDAH" character varying(100),
    "UNIT_ASAL" character varying(32),
    "UNIT_TUJUAN" character varying(32),
    "SURAT_PERNYATAAN_MELEPAS" character varying(100),
    "SK_KP_TERAKHIR" character varying(100),
    "SK_JABATAN" character varying(100),
    "SKP" character varying(10),
    "SK_TUNKIN" character varying(100),
    "SURAT_PERNYATAAN_MENERIMA" character varying(100),
    "NO_SK_PINDAH" character varying(100),
    "TANGGAL_SK_PINDAH" character varying(10),
    "FILE_SK" character varying(100),
    "STATUS_SATKER" integer,
    "STATUS_BIRO" integer,
    "JABATAN_ID" numeric,
    "KETERANGAN" character(255),
    "TANGGAL_TMT_PINDAH" date,
    "CREATED_DATE" date,
    "CREATED_BY" integer
);

CREATE SEQUENCE "pindah_unit_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "pindah_unit_ID_seq" OWNED BY pindah_unit."ID";

CREATE VIEW pns_aktif_old AS
 SELECT "ID",
    masa_kerja[0] AS masa_kerja_th,
    masa_kerja[1] AS masa_kerja_bl
   FROM ( SELECT pegawai."ID",
            get_masa_kerja_arr(pegawai."TMT_CPNS", ('now'::text)::date) AS masa_kerja
           FROM pegawai
          WHERE ((pegawai.status_pegawai = 1) AND ((pegawai.terminated_date IS NULL) OR ((pegawai.terminated_date IS NOT NULL) AND (pegawai.terminated_date > ('now'::text)::date))))) temp;

CREATE TABLE ref_jabatan (
    "ID_JABATAN" double precision NOT NULL,
    "NAMA_JABATAN" text,
    "JENIS_JABATAN" character varying(100),
    "KELAS" smallint,
    "PENSIUN" smallint,
    "KODE_BKN" character varying(250),
    "TUNJANGAN" double precision
);

CREATE TABLE ref_tunjangan_jabatan (
    "ID_TUNJAB" integer NOT NULL,
    "ESELON" character varying(10),
    "BESARAN_TUNJAB" character varying(100)
);

CREATE TABLE ref_tunjangan_kinerja (
    "ID" integer NOT NULL,
    "KELAS_JABATAN" integer,
    "TUNJANGAN_KINERJA" double precision
);

CREATE SEQUENCE "ref_tunjangan_kinerja_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "ref_tunjangan_kinerja_ID_seq" OWNED BY ref_tunjangan_kinerja."ID";

CREATE VIEW rekap_agama_jenis_kelamin AS
 SELECT pegawai."JENIS_KELAMIN",
    agama."ID",
    agama."NAMA",
    count(*) AS total
   FROM (pegawai
     LEFT JOIN agama ON ((pegawai."AGAMA_ID" = agama."ID")))
  GROUP BY pegawai."JENIS_KELAMIN", agama."ID", agama."NAMA"
  ORDER BY agama."NAMA";

CREATE TABLE request_formasi (
    id bigint NOT NULL,
    unit_id character varying(32),
    jumlah_ajuan smallint,
    kualifikasi_pendidikan text,
    id_jabatan character varying(32),
    satker_id character varying(32),
    tahun character varying(4),
    skala_prioritas smallint
);

CREATE SEQUENCE request_formasi_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE request_formasi_id_seq OWNED BY request_formasi.id;

CREATE TABLE role_permissions (
    role_id bigint,
    permission_id bigint,
    id bigint NOT NULL
);

CREATE SEQUENCE role_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE role_permissions_id_seq OWNED BY role_permissions.id;

CREATE SEQUENCE roles_role_id_seq
    START WITH 6
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE roles (
    role_name character(60) NOT NULL,
    description character varying(255) DEFAULT NULL::character varying,
    is_default integer DEFAULT 0 NOT NULL,
    can_delete integer DEFAULT 1 NOT NULL,
    login_destination character(255) DEFAULT '/'::bpchar NOT NULL,
    deleted integer DEFAULT 0 NOT NULL,
    default_context character(255) DEFAULT 'content'::bpchar NOT NULL,
    role_id integer DEFAULT nextval('roles_role_id_seq'::regclass) NOT NULL
);

CREATE TABLE roles_users (
    role_id bigint,
    user_id bigint,
    role_user_id bigint NOT NULL
);

CREATE SEQUENCE roles_users_role_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE roles_users_role_user_id_seq OWNED BY roles_users.role_user_id;

CREATE TABLE rpt_golongan_bulan (
    "ID" smallint NOT NULL,
    "GOLONGAN_ID" character varying(10),
    "GOLONGAN_NAMA" character varying(40),
    "BULAN" smallint,
    "TAHUN" smallint,
    "JUMLAH" smallint
);

CREATE SEQUENCE "rpt_golongan_bulan_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rpt_golongan_bulan_ID_seq" OWNED BY rpt_golongan_bulan."ID";

CREATE TABLE rpt_jumlah_asn (
    "ID" integer NOT NULL,
    "BULAN" smallint,
    "TAHUN" character varying(4),
    "JENIS" character varying(20),
    "KETERANGAN" character varying(50),
    "JUMLAH" real
);

CREATE SEQUENCE "rpt_jumlah_asn_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rpt_jumlah_asn_ID_seq" OWNED BY rpt_jumlah_asn."ID";

CREATE SEQUENCE rpt_jumlah_asn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rpt_jumlah_asn_id_seq OWNED BY rpt_jumlah_asn."ID";

CREATE TABLE rpt_pendidikan_bulan (
    "ID" smallint NOT NULL,
    "TINGKAT_PENDIDIKAN" character varying(16),
    "NAMA_TINGKAT" character varying(255),
    "BULAN" smallint,
    "TAHUN" smallint,
    "JUMLAH" smallint
);

CREATE SEQUENCE "rpt_pendidikan_bulan_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rpt_pendidikan_bulan_ID_seq" OWNED BY rpt_pendidikan_bulan."ID";

CREATE SEQUENCE "rwt_assesmen_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rwt_assesmen_ID_seq" OWNED BY rwt_assesmen."ID";

CREATE TABLE rwt_diklat (
    id bigint NOT NULL,
    jenis_diklat character varying,
    jenis_diklat_id character varying,
    institusi_penyelenggara character varying,
    nomor_sertifikat character varying,
    tanggal_mulai date,
    tanggal_selesai date,
    tahun_diklat integer,
    durasi_jam integer,
    pns_orang_id character varying,
    nip_baru character varying,
    createddate timestamp without time zone DEFAULT now(),
    diklat_struktural_id character varying,
    nama_diklat character varying,
    file_base64 text,
    rumpun_diklat character varying,
    rumpun_diklat_id character varying,
    sudah_kirim_siasn character varying DEFAULT 'belum'::character varying,
    siasn_id character varying
);

CREATE SEQUENCE "rwt_diklat_ID_seq"
    START WITH 37
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE "rwt_diklat_fungsional_ID_seq"
    START WITH 7
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_diklat_fungsional (
    "DIKLAT_FUNGSIONAL_ID" character varying(255) DEFAULT nextval('"rwt_diklat_fungsional_ID_seq"'::regclass) NOT NULL,
    "NIP_BARU" character varying(255),
    "NIP_LAMA" character varying(255),
    "JENIS_DIKLAT" character varying(255),
    "NAMA_KURSUS" character varying(255),
    "JUMLAH_JAM" character varying(255),
    "TAHUN" character varying(255),
    "INSTITUSI_PENYELENGGARA" character varying(255),
    "JENIS_KURSUS_SERTIPIKAT" character varying(255),
    "NOMOR_SERTIPIKAT" character varying(255),
    "INSTANSI" character varying(255),
    "STATUS_DATA" character varying(255),
    "TANGGAL_KURSUS" date,
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(200),
    "LAMA" real,
    "SIASN_ID" character varying
);

CREATE SEQUENCE rwt_diklat_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_diklat_id_seq OWNED BY rwt_diklat.id;

CREATE TABLE rwt_diklat_struktural (
    "ID" character varying(255) DEFAULT nextval('"rwt_diklat_ID_seq"'::regclass) NOT NULL,
    "PNS_ID" character varying(255),
    "PNS_NIP" character varying(255),
    "PNS_NAMA" character varying(255),
    "ID_DIKLAT" character varying(255),
    "NAMA_DIKLAT" character varying(255),
    "NOMOR" character varying(255),
    "TANGGAL" date,
    "TAHUN" smallint,
    "STATUS_DATA" character varying(15),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(200),
    "LAMA" real,
    "CREATED_DATE" timestamp without time zone DEFAULT now(),
    "SIASN_ID" character varying
);

CREATE SEQUENCE rwt_golongan_id_seq
    START WITH 351278
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_golongan (
    "ID" character varying(255) DEFAULT nextval('rwt_golongan_id_seq'::regclass) NOT NULL,
    "PNS_ID" character varying(255),
    "PNS_NIP" character varying(255),
    "PNS_NAMA" character varying(255),
    "KODE_JENIS_KP" character varying(255),
    "JENIS_KP" character varying(255),
    "ID_GOLONGAN" character varying(255),
    "GOLONGAN" character varying(255),
    "PANGKAT" character varying(255),
    "SK_NOMOR" character varying(255),
    "NOMOR_BKN" character varying(255),
    "JUMLAH_ANGKA_KREDIT_UTAMA" character varying(255),
    "JUMLAH_ANGKA_KREDIT_TAMBAHAN" character varying(255),
    "MK_GOLONGAN_TAHUN" character varying(255),
    "MK_GOLONGAN_BULAN" character varying(255),
    "SK_TANGGAL" date,
    "TANGGAL_BKN" date,
    "TMT_GOLONGAN" date,
    "STATUS_SATKER" integer,
    "STATUS_BIRO" integer,
    "PANGKAT_TERAKHIR" integer,
    "ID_BKN" character varying(36),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255),
    id_arsip bigint,
    "GOLONGAN_ASAL" character varying(2),
    "BASIC" character varying(15),
    "SK_TYPE" smallint,
    "KANREG" character varying(5),
    "KPKN" character varying(50),
    "KETERANGAN" character varying(255),
    "LPNK" character varying(10),
    "JENIS_RIWAYAT" character varying(50)
);

CREATE TABLE rwt_hukdis (
    "ID" bigint NOT NULL,
    "PNS_ID" character(32),
    "PNS_NIP" character(21),
    "NAMA" character(200),
    "ID_GOLONGAN" character(2),
    "NAMA_GOLONGAN" character(20),
    "ID_JENIS_HUKUMAN" character(2),
    "NAMA_JENIS_HUKUMAN" character(100),
    "SK_NOMOR" character(30),
    "SK_TANGGAL" date,
    "TANGGAL_MULAI_HUKUMAN" date,
    "MASA_TAHUN" integer,
    "MASA_BULAN" integer,
    "TANGGAL_AKHIR_HUKUMAN" date,
    "NO_PP" character(20),
    "NO_SK_PEMBATALAN" character(20),
    "TANGGAL_SK_PEMBATALAN" date,
    "ID_BKN" character varying(255),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255)
);

CREATE SEQUENCE "rwt_hukdis_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rwt_hukdis_ID_seq" OWNED BY rwt_hukdis."ID";

CREATE SEQUENCE "rwt_jabatan_ID_seq"
    START WITH 94276
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_jabatan (
    "ID_BKN" character(64),
    "PNS_ID" character(100),
    "PNS_NIP" character(25),
    "PNS_NAMA" character(200),
    "ID_UNOR" character(100),
    "UNOR" text,
    "ID_JENIS_JABATAN" character(10),
    "JENIS_JABATAN" character(250),
    "ID_JABATAN" character(100),
    "NAMA_JABATAN" text,
    "ID_ESELON" character(32),
    "ESELON" character(100),
    "TMT_JABATAN" date,
    "NOMOR_SK" character(100),
    "TANGGAL_SK" date,
    "ID_SATUAN_KERJA" character varying(36),
    "TMT_PELANTIKAN" date,
    "IS_ACTIVE" character(1),
    "ESELON1" text,
    "ESELON2" text,
    "ESELON3" text,
    "ESELON4" text,
    "ID" bigint DEFAULT nextval('"rwt_jabatan_ID_seq"'::regclass) NOT NULL,
    "CATATAN" character(255),
    "JENIS_SK" character(100),
    "LAST_UPDATED" date,
    "STATUS_SATKER" integer,
    "STATUS_BIRO" integer,
    "ID_JABATAN_BKN" character varying(36),
    "ID_UNOR_BKN" character varying(36),
    "JABATAN_TERAKHIR" integer,
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255),
    "ID_TABEL_MUTASI" bigint,
    "TERMINATED_DATE" date
);

CREATE TABLE rwt_jabatan_empty (
    "ID_BKN" character(64),
    "PNS_ID" character(100),
    "PNS_NIP" character(25),
    "PNS_NAMA" character(200),
    "ID_UNOR" character(100),
    "UNOR" text,
    "ID_JENIS_JABATAN" character(10),
    "JENIS_JABATAN" character(250),
    "ID_JABATAN" character(100),
    "NAMA_JABATAN" text,
    "ID_ESELON" character(32),
    "ESELON" character(100),
    "TMT_JABATAN" date,
    "NOMOR_SK" character(100),
    "TANGGAL_SK" date,
    "ID_SATUAN_KERJA" character(32),
    "TMT_PELANTIKAN" date,
    "IS_ACTIVE" character(1),
    "ESELON1" text,
    "ESELON2" text,
    "ESELON3" text,
    "ESELON4" text,
    "ID" bigint DEFAULT nextval('"rwt_jabatan_ID_seq"'::regclass) NOT NULL,
    "CATATAN" character(255),
    "JENIS_SK" character(100),
    "LAST_UPDATED" date,
    "STATUS_SATKER" integer,
    "STATUS_BIRO" integer,
    "ID_JABATAN_BKN" character(64),
    "ID_UNOR_BKN" character(32),
    "JABATAN_TERAKHIR" integer,
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255),
    "ID_TABEL_MUTASI" bigint,
    "TERMINATED_DATE" date
);

CREATE TABLE rwt_kgb (
    pegawai_id integer,
    tmt_sk date,
    alasan character varying(255),
    mv_kgb_id bigint,
    no_sk character varying(255),
    pejabat character varying(255),
    id bigint NOT NULL,
    ref character varying(255) DEFAULT uuid_generate_v4(),
    tgl_sk date,
    pegawai_nama character varying(255),
    pegawai_nip character varying(255),
    birth_place character varying(255),
    birth_date date,
    o_gol_ruang character varying(255),
    o_gol_tmt character varying(255),
    o_masakerja_thn smallint,
    o_masakerja_bln smallint,
    o_gapok character varying(255),
    o_jabatan_text character varying(255),
    o_tmt_jabatan date,
    n_gol_ruang character varying(255),
    n_gol_tmt character varying(255),
    n_masakerja_thn smallint,
    n_masakerja_bln smallint,
    n_gapok character varying(255),
    n_jabatan_text character varying(255),
    n_tmt_jabatan date,
    n_golongan_id integer,
    unit_kerja_text character varying(255),
    unit_kerja_induk_text character varying(255),
    unit_kerja_induk_id character varying(255),
    kantor_pembayaran character varying(255),
    last_education character varying(255),
    last_education_date date,
    nama_pejabat character varying(100),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255)
);

CREATE SEQUENCE rwt_kgb_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_kgb_id_seq OWNED BY rwt_kgb.id;

CREATE TABLE rwt_kinerja (
    id integer NOT NULL,
    id_simarin integer,
    tahun integer,
    periode_mulai date,
    periode_selesai date,
    format_skp character varying(20),
    jenis_skp character varying(20),
    idp_pegawai character varying(7),
    nip character varying(20),
    nama character varying(200),
    panggol character varying(50),
    jabatan character varying(200),
    penugasan character varying(200),
    id_unit_kerja integer,
    unit_kerja character varying(200),
    idp_penilai character varying(7),
    nip_penilai character varying(20),
    nama_penilai character varying(200),
    panggol_penilai character varying(50),
    jabatan_penilai character varying(200),
    penugasan_penilai character varying(200),
    id_unit_kerja_penilai integer,
    unit_kerja_penilai character varying(200),
    idp_atasan_penilai character varying(7),
    nip_atasan_penilai character varying(20),
    nama_atasan_penilai character varying(200),
    panggol_atasan_penilai character varying(50),
    jabatan_atasan_penilai character varying(200),
    penugasan_atasan_penilai character varying(200),
    id_unit_kerja_atasan_penilai integer,
    unit_kerja_atasan_penilai character varying(200),
    idp_penilai_simarin character varying(7),
    nip_penilai_simarin character varying(20),
    nama_penilai_simarin character varying(200),
    panggol_penilai_simarin character varying(50),
    jabatan_penilai_simarin character varying(200),
    penugasan_penilai_simarin character varying(200),
    id_unit_kerja_penilai_simarin integer,
    unit_kerja_penilai_simarin character varying(200),
    idp_penilai_realisasi character varying(7),
    nip_penilai_realisasi character varying(20),
    nama_penilai_realisasi character varying(200),
    panggol_penilai_realisasi character varying(50),
    jabatan_penilai_realisasi character varying(200),
    penugasan_penilai_realisasi character varying(200),
    id_unit_kerja_penilai_realisasi integer,
    unit_kerja_penilai_realisasi character varying(200),
    idp_atasan_penilai_realisasi character varying(7),
    nip_atasan_penilai_realisasi character varying(20),
    nama_atasan_penilai_realisasi character varying(200),
    panggol_atasan_penilai_realisasi character varying(50),
    jabatan_atasan_penilai_realisasi character varying(200),
    penugasan_atasan_penilai_realisasi character varying(200),
    id_unit_kerja_atasan_penilai_realisasi integer,
    unit_kerja_atasan_penilai_realisasi character varying(200),
    idp_penilai_realisasi_simarin character varying(7),
    nip_penilai_realisasi_simarin character varying(20),
    nama_penilai_realisasi_simarin character varying(200),
    panggol_penilai_realisasi_simarin character varying(50),
    jabatan_penilai_realisasi_simarin character varying(200),
    penugasan_penilai_realisasi_simarin character varying(200),
    id_unit_kerja_penilai_realisasi_simarin integer,
    unit_kerja_penilai_realisasi_simarin character varying(200),
    nama_realisasi character varying(200),
    nip_realisasi character varying(20),
    panggol_realisasi character varying(50),
    jabatan_realisasi character varying(200),
    penugasan_realisasi character varying(200),
    id_unit_kerja_realisasi integer,
    unit_kerja_realisasi character varying(200),
    skp_instansi_lama character varying(200),
    capaian_kinerja_org character varying(20),
    pola_distribusi_img character varying(200),
    nilai_akhir_hasil_kerja real,
    rating_hasil_kerja character varying(50),
    nilai_akhir_perilaku_kerja real,
    rating_perilaku_kerja character varying(50),
    predikat_kinerja character varying(100),
    tunjangan_kinerja integer,
    catatan_rekomendasi text,
    is_keberatan character varying(100),
    keberatan text,
    penjelasan_pejabat_penilai text,
    keputusan_rekomendasi_atasan_pejabat text,
    url_skp_instansi_lama character varying(200),
    is_keberatan_date character varying(5),
    ref uuid DEFAULT uuid_generate_v4(),
    id_arsip integer,
    created_date timestamp(6) without time zone DEFAULT now()
);

CREATE SEQUENCE rwt_kinerja_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_kinerja_id_seq OWNED BY rwt_kinerja.id;

CREATE TABLE rwt_kursus (
    "PNS_ID_x" character(32),
    "PNS_NIP" character(30),
    "TIPE_KURSUS" character(10),
    "JENIS_KURSUS" character(30),
    "NAMA_KURSUS" character(200),
    "LAMA_KURSUS" double precision,
    "TANGGAL_KURSUS" date,
    "NO_SERTIFIKAT" character(50),
    "INSTANSI" character(200),
    "INSTITUSI_PENYELENGGARA" character(200),
    "ID" integer NOT NULL,
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(200),
    "CREATEDDATE" timestamp without time zone DEFAULT now(),
    "SIASN_ID" character varying,
    "PNS_ID" character varying
);

CREATE SEQUENCE "rwt_kursus_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rwt_kursus_ID_seq" OWNED BY rwt_kursus."ID";

CREATE SEQUENCE "rwt_pekerjaan_ID_seq"
    START WITH 3
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_pekerjaan (
    "ID" integer DEFAULT nextval('"rwt_pekerjaan_ID_seq"'::regclass) NOT NULL,
    "PNS_NIP" character(30),
    "JENIS_PERUSAHAAN" character(100),
    "NAMA_PERUSAHAAN" character(200),
    "SEBAGAI" character(200),
    "DARI_TANGGAL" date,
    "SAMPAI_TANGGAL" date,
    "PNS_ID" character(32),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255)
);

CREATE SEQUENCE "rwt_pendidikan_ID_seq"
    START WITH 66872
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_pendidikan (
    "ID" integer DEFAULT nextval('"rwt_pendidikan_ID_seq"'::regclass) NOT NULL,
    "PNS_ID_3" character varying(32),
    "TINGKAT_PENDIDIKAN_ID" character varying(32),
    "PENDIDIKAN_ID_3" character varying(32),
    "TANGGAL_LULUS" date,
    "NOMOR_IJASAH" character varying(100),
    "NAMA_SEKOLAH" character varying(200),
    "GELAR_DEPAN" character varying(50),
    "GELAR_BELAKANG" character varying(60),
    "PENDIDIKAN_PERTAMA" character varying(1),
    "NEGARA_SEKOLAH" character varying(255),
    "TAHUN_LULUS" character varying(4),
    "NIP" character(35),
    "DIAKUI_BKN" integer,
    "TUGAS_BELAJAR" character(255),
    "STATUS_SATKER" integer,
    "STATUS_BIRO" integer,
    "PENDIDIKAN_TERAKHIR" integer,
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(200),
    "PNS_ID" character varying(255) NOT NULL,
    "PENDIDIKAN_ID" character varying
);

CREATE SEQUENCE "rwt_penghargaan_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_penghargaan (
    "ID" integer DEFAULT nextval('"rwt_penghargaan_ID_seq"'::regclass) NOT NULL,
    "PNS_ID" character(32),
    "PNS_NIP" character(21),
    "NAMA" character(200),
    "ID_GOLONGAN" character(2),
    "NAMA_GOLONGAN" character(100),
    "ID_JENIS_PENGHARGAAN" character(3),
    "NAMA_JENIS_PENGHARGAAN" character(100),
    "SK_NOMOR" character(30),
    "SK_TANGGAL" date,
    "ID_BKN" character varying(255),
    "SURAT_USUL" text,
    "KETERANGAN" text
);

CREATE TABLE rwt_penghargaan_umum (
    id bigint NOT NULL,
    jenis_penghargaan character varying,
    deskripsi_penghargaan character varying,
    tanggal_penghargaan date,
    createddate timestamp without time zone DEFAULT now(),
    exist boolean DEFAULT true
);

CREATE SEQUENCE rwt_penghargaan_umum_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_penghargaan_umum_id_seq OWNED BY rwt_penghargaan_umum.id;

CREATE TABLE rwt_penugasan (
    id bigint NOT NULL,
    tipe_jabatan character varying,
    deskripsi_jabatan text,
    tanggal_mulai date,
    tanggal_selesai date,
    createddate timestamp without time zone DEFAULT now(),
    exist boolean DEFAULT true
);

CREATE SEQUENCE rwt_penugasan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_penugasan_id_seq OWNED BY rwt_penugasan.id;

CREATE SEQUENCE rwt_pindah_unit_kerja_id_seq
    START WITH 2
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_pindah_unit_kerja (
    "ID" character varying(255) DEFAULT nextval('rwt_pindah_unit_kerja_id_seq'::regclass) NOT NULL,
    "PNS_ID" character varying(255),
    "PNS_NIP" character varying(255),
    "PNS_NAMA" character varying(255),
    "SK_NOMOR" character varying(255),
    "ASAL_ID" character varying(255),
    "ASAL_NAMA" character varying(255),
    "ID_UNOR_BARU" character varying(255),
    "NAMA_UNOR_BARU" character varying(255),
    "ID_INSTANSI" character varying(255),
    "NAMA_INSTANSI" character varying(255),
    "SK_TANGGAL" date,
    "ID_SATUAN_KERJA" character varying(100),
    "NAMA_SATUAN_KERJA" character varying(255),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255)
);

CREATE TABLE rwt_pns_cpns (
    "ID" integer NOT NULL,
    "STATUS_KEPEGAWAIAN" character varying(5),
    "TMT_CPNS" date,
    "TGL_SK_CPNS" date,
    "NO_SK_CPNS" character varying(100),
    "JENIS_PENGADAAN" character varying(30),
    "TGL_SPMT" date,
    "NO_SPMT" character varying(100),
    "TMT_PNS" date,
    "TGL_SK_PNS" date,
    "N0_SK_PNS" character varying(100),
    "TGL_PERTEK_C2TH" date,
    "NO_PERTEK_C2TH" character varying(100),
    "TGL_KEP_HONORER_2TAHUN" date,
    "NO_PERTEK_KEP_HONORER_2TAHUN" character varying(50),
    "KARIS_KARSU" character varying(20),
    "KARPEG" character varying(20),
    "TGL_STTPL" date,
    "NO_STTPL" character varying(100),
    "TGL_DOKTER" date,
    "NO_SURAT_DOKTER" character varying(50),
    "NAMA_JABATAN_ANGKAT_CPNS" character varying(200),
    "PNS_ID" character varying(32),
    "PNS_NIP" character varying(18)
);

CREATE SEQUENCE "rwt_pns_cpns_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rwt_pns_cpns_ID_seq" OWNED BY rwt_pns_cpns."ID";

CREATE SEQUENCE rwt_prestasi_kerja_id_seq
    START WITH 16
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE rwt_prestasi_kerja (
    "ID" character varying(255) DEFAULT nextval('rwt_prestasi_kerja_id_seq'::regclass) NOT NULL,
    "PNS_NIP" character varying(255),
    "PNS_NAMA" character varying(255),
    "ATASAN_LANGSUNG_PNS_NAMA" character varying(255),
    "ATASAN_LANGSUNG_PNS_NIP" character varying(255),
    "NILAI_SKP" character varying(255),
    "NILAI_PROSENTASE_SKP" character varying(255),
    "NILAI_SKP_AKHIR" character varying(255),
    "PERILAKU_KOMITMEN" character varying(255),
    "PERILAKU_INTEGRITAS" character varying(255),
    "PERILAKU_DISIPLIN" character varying(255),
    "PERILAKU_KERJASAMA" character varying(255),
    "PERILAKU_ORIENTASI_PELAYANAN" character varying(255),
    "PERILAKU_KEPEMIMPINAN" character varying(255),
    "NILAI_PERILAKU" character varying(255),
    "NILAI_PROSENTASE_PERILAKU" character varying(255),
    "NILAI_PERILAKU_AKHIR" character varying(255),
    "NILAI_PPK" character varying(255),
    "TAHUN" integer,
    "JABATAN_TIPE" character varying(255),
    "PNS_ID" character varying(255),
    "ATASAN_LANGSUNG_PNS_ID" character varying(255),
    "ATASAN_ATASAN_LANGSUNG_PNS_ID" character varying(255),
    "ATASAN_ATASAN_LANGSUNG_PNS_NAMA" character varying(255),
    "ATASAN_ATASAN_LANGSUNG_PNS_NIP" character varying(255),
    "JABATAN_TIPE_TEXT" character varying(255),
    "ATASAN_LANGSUNG_PNS_JABATAN" character varying(255),
    "ATASAN_ATASAN_LANGSUNG_PNS_JABATAN" character varying(255),
    "JABATAN_NAMA" character varying(255),
    "BKN_ID" character varying(36),
    "UNOR_PENILAI" character varying(200),
    "UNOR_ATASAN_PENILAI" character varying(200),
    "ATASAN_PENILAI_PNS" character varying(200),
    "PENILAI_PNS" character varying(200),
    "GOL_PENILAI" character varying(20),
    "GOL_ATASAN_PENILAI" character varying(20),
    "TMT_GOL_PENILAI" character varying(20),
    "TMT_GOL_ATASAN_PENILAI" character varying(255),
    "PERATURAN" character varying(20),
    created_date date,
    updated_date date,
    "PERILAKU_INISIATIF_KERJA" character varying(20)
);

CREATE TABLE rwt_tugas_belajar (
    "ID" integer NOT NULL,
    "NAMA" character varying(100),
    "NIP" character varying(30),
    "TINGKAT_PENDIDIKAN" character varying(5),
    "PROGRAM_STUDI" character varying(200),
    "FAKULTAS" character varying(100),
    "UNIVERSITAS" character varying(100),
    "MULAI_BELAJAR" date,
    "AKHIR_BELAJAR" date,
    "NOMOR_SK" character varying(50),
    "TANGGAL_SK" date,
    "KETERANGAN" character varying(200),
    "JENIS_USUL" character varying(20),
    "ID_PENGAJUAN" character varying(20),
    "FILE_BASE64" text,
    "KETERANGAN_BERKAS" character varying(255)
);

CREATE SEQUENCE "rwt_tugas_belajar_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "rwt_tugas_belajar_ID_seq" OWNED BY rwt_tugas_belajar."ID";

CREATE TABLE rwt_ujikom (
    id bigint NOT NULL,
    jenis_ujikom character varying,
    nip_baru character varying,
    link_sertifikat character varying,
    createddate timestamp without time zone DEFAULT now(),
    exist boolean,
    tahun integer
);

CREATE SEQUENCE rwt_ujikom_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE rwt_ujikom_id_seq OWNED BY rwt_ujikom.id;

CREATE TABLE settings (
    name text NOT NULL,
    module character varying(50) NOT NULL,
    value character varying(500) NOT NULL,
    id bigint NOT NULL
);

CREATE SEQUENCE settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE settings_id_seq OWNED BY settings.id;

CREATE TABLE sisa_cuti (
    "ID" integer NOT NULL,
    "PNS_NIP" character varying(18) NOT NULL,
    "TAHUN" character varying(4) NOT NULL,
    "SISA_N" smallint,
    "SISA_N_1" smallint,
    "SISA_N_2" smallint,
    "SISA" smallint,
    "NAMA" character varying(100),
    "SUDAH_DIAMBIL" smallint
);

CREATE SEQUENCE "sisa_cuti_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "sisa_cuti_ID_seq" OWNED BY sisa_cuti."ID";

CREATE TABLE synch_jumlah_pegawai (
    id bigint NOT NULL,
    kode_unit_kerja character varying(20),
    id_unor_bkn character varying(32),
    nama_eselon_1 character varying(200),
    satker character varying(200),
    satker_singkatan character varying(100),
    jumlah_mutasi integer,
    jumlah_dikbudhr integer,
    update_time timestamp(6) without time zone
);

CREATE SEQUENCE synch_jumlah_pegawai_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE synch_jumlah_pegawai_id_seq OWNED BY synch_jumlah_pegawai.id;

CREATE TABLE tb_nomor_batasan (
    "BATAS_AWAL" character varying(100) NOT NULL,
    "BATAS_AKHIR" character varying(100) NOT NULL,
    "TAHUN_NOMOR" character varying(100)
);

CREATE TABLE tb_nomor_surat (
    id double precision NOT NULL,
    nomor_surat double precision NOT NULL,
    kode character varying(100) NOT NULL,
    tanggal date NOT NULL,
    kepada text NOT NULL,
    keterangan text NOT NULL,
    username character varying(100) NOT NULL
);

CREATE TABLE tbl_cek (
    id_file character varying(200)
);

CREATE TABLE tbl_file_ds_corrector (
    korektor_ke smallint,
    id_pegawai_korektor character varying(100),
    is_returned smallint,
    catatan_koreksi text,
    is_corrected smallint,
    id_file character varying(200),
    id integer NOT NULL
);

COMMENT ON COLUMN tbl_file_ds_corrector.is_returned IS '1=dikembalikan, 0/null = sudah oke';

COMMENT ON COLUMN tbl_file_ds_corrector.is_corrected IS '1=koreksi ok, 2=siap koreksi, 0/null = masih antrian';

CREATE SEQUENCE tbl_file_ds_corrector_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tbl_file_ds_corrector_id_seq OWNED BY tbl_file_ds_corrector.id;

CREATE SEQUENCE tbl_file_ds_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tbl_file_ds_id_seq OWNED BY tbl_file_ds.id;

CREATE TABLE tbl_file_ds_khusus_login (
    "ID_FILE" character varying(255) NOT NULL
);

CREATE TABLE tbl_file_ds_riwayat (
    id_file character varying(200),
    id_pemroses character varying(255),
    tindakan text,
    catatan_tindakan text,
    waktu_tindakan timestamp(6) without time zone,
    akses_pengguna character varying(200),
    id_riwayat bigint NOT NULL
);

CREATE SEQUENCE tbl_file_ds_riwayat_id_riwayat_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tbl_file_ds_riwayat_id_riwayat_seq OWNED BY tbl_file_ds_riwayat.id_riwayat;

CREATE TABLE tbl_file_ttd (
    id_pns_bkn character varying(200) NOT NULL,
    nip character varying(50),
    base64ttd text
);

CREATE TABLE tbl_kategori_dokumen (
    id_kategori smallint NOT NULL,
    kategori_dokumen character varying(255),
    update_jabatan smallint,
    via_ds smallint,
    kelompok character varying(255),
    kaitan character varying(255),
    grup_proses character varying(255),
    izinkan_kolektif character varying(10),
    grup_info character varying(200),
    untuk_pegawai smallint,
    login_untuk_lihat smallint
);

CREATE TABLE tbl_kategori_dokumen_penandatangan (
    "ID_URUT" bigint NOT NULL,
    "KELOMPOK" character varying(255),
    "PENANDATANGAN" character varying(255),
    "KOREKTOR_KE" smallint,
    "NAMA_KOREKTOR" character varying(255),
    "JABATAN" text,
    "SATKER" text,
    "ID_PNS" character varying(255),
    "ID_UNOR" character varying(255)
);

CREATE TABLE tbl_pengantar_dokumen (
    id_pengantar character varying(255) NOT NULL,
    html_lampiran text,
    skema character varying(10)
);

CREATE TABLE tte_master_variable (
    id smallint NOT NULL,
    label_variable character varying(50),
    nama_variable character varying(50),
    tipe character varying(50),
    keterangan character varying(255)
);

CREATE SEQUENCE "tte_ master_variable_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "tte_ master_variable_id_seq" OWNED BY tte_master_variable.id;

CREATE TABLE tte_master_korektor (
    id smallint NOT NULL,
    id_tte_master_proses smallint,
    id_pegawai_korektor character varying(32),
    korektor_ke smallint
);

CREATE SEQUENCE tte_master_korektor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_master_korektor_id_seq OWNED BY tte_master_korektor.id;

CREATE TABLE tte_master_proses (
    id integer NOT NULL,
    nama_proses character varying(100) NOT NULL,
    template_sk character varying(100),
    penandatangan_sk character varying(32),
    keterangan_proses text
);

CREATE SEQUENCE tte_master_proses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_master_proses_id_seq OWNED BY tte_master_proses.id;

CREATE TABLE tte_master_proses_variable (
    id smallint NOT NULL,
    id_proses smallint,
    id_variable smallint
);

CREATE SEQUENCE tte_master_proses_variable_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_master_proses_variable_id_seq OWNED BY tte_master_proses_variable.id;

CREATE TABLE tte_trx_draft_sk (
    id bigint NOT NULL,
    id_master_proses smallint,
    nip_sk character varying(32),
    penandatangan_sk character varying(32),
    tgl_sk date,
    nomor_sk character varying(100),
    file_template character varying(100),
    base64pdf_hasil text,
    created_date date,
    created_by bigint,
    updated_date date,
    updated_by bigint,
    id_file character varying(40),
    tmt_sk date,
    nama_pemilik_sk character varying(255),
    halaman_ttd smallint DEFAULT 1,
    show_qrcode smallint DEFAULT 0,
    letak_ttd smallint DEFAULT 0
);

COMMENT ON COLUMN tte_trx_draft_sk.nama_pemilik_sk IS 'nama pemilik sk';

CREATE TABLE tte_trx_draft_sk_detil (
    id integer NOT NULL,
    id_tte_trx_draft_sk bigint,
    id_variable smallint,
    isi character varying(255)
);

CREATE SEQUENCE tte_trx_draft_sk_detil_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_trx_draft_sk_detil_id_seq OWNED BY tte_trx_draft_sk_detil.id;

CREATE SEQUENCE tte_trx_draft_sk_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_trx_draft_sk_id_seq OWNED BY tte_trx_draft_sk.id;

CREATE TABLE tte_trx_korektor_draft (
    id integer NOT NULL,
    id_tte_trx_draft_sk bigint,
    id_pegawai_korektor character varying(32),
    korektor_ke smallint
);

CREATE SEQUENCE tte_trx_korektor_draft_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tte_trx_korektor_draft_id_seq OWNED BY tte_trx_korektor_draft.id;

CREATE SEQUENCE unitkerja_id_seq
    START WITH 1245
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE unitkerja_old_id_seq
    START WITH 1243
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE update_mandiri (
    "ID" integer NOT NULL,
    "PNS_ID" character(32),
    "KOLOM" character(70),
    "DARI" character(400),
    "PERUBAHAN" character(400),
    "STATUS" integer,
    "VERIFIKASI_BY" integer,
    "VERIFIKASI_TGL" date,
    "UPDATE_TGL" date,
    "NAMA_KOLOM" character(100),
    "LEVEL_UPDATE" integer,
    "ID_TABEL" integer,
    "UPDATED_BY" integer
);

CREATE SEQUENCE "update_mandiri_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE "update_mandiri_ID_seq" OWNED BY update_mandiri."ID";

CREATE SEQUENCE user_id_seq
    START WITH 32952
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE user_meta_meta_id_seq
    START WITH 414
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE users (
    id bigint DEFAULT nextval('user_id_seq'::regclass) NOT NULL,
    role_id bigint DEFAULT (4)::bigint NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(255) DEFAULT ''::bpchar NOT NULL,
    password_hash character varying(255) NOT NULL,
    reset_hash character varying(255) DEFAULT NULL::bpchar,
    last_login date,
    last_ip character(40) DEFAULT ''::bpchar NOT NULL,
    created_on date,
    deleted integer DEFAULT 0 NOT NULL,
    reset_by integer,
    banned integer DEFAULT 0 NOT NULL,
    ban_message character(255) DEFAULT NULL::bpchar,
    display_name character(255) DEFAULT ''::bpchar,
    display_name_changed date,
    timezone character(4) DEFAULT 'UM7'::bpchar NOT NULL,
    language character(20) DEFAULT 'english'::bpchar NOT NULL,
    active integer DEFAULT 0 NOT NULL,
    activate_hash character(40) DEFAULT ''::bpchar NOT NULL,
    password_iterations integer NOT NULL,
    force_password_reset integer DEFAULT 0 NOT NULL,
    nip character varying(20) DEFAULT NULL::bpchar,
    satkers text,
    admin_nomor smallint,
    imei character varying(100),
    token character varying(255),
    real_imei character varying(100),
    fcm character varying(255),
    banned_asigo integer DEFAULT 0
);

CREATE TABLE usulan_dokumen (
    id bigint NOT NULL,
    perkiraan_id bigint,
    title character varying(255),
    file_upload character varying(255),
    _created_at timestamp(6) without time zone DEFAULT now(),
    _created_by bigint,
    tipe character varying(255)
);

CREATE SEQUENCE usulan_documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE usulan_documents_id_seq OWNED BY usulan_dokumen.id;

CREATE VIEW v_kategori_ds AS
 SELECT DISTINCT kategori AS kategori_ds
   FROM tbl_file_ds
  ORDER BY kategori;

CREATE VIEW vw_unor_satker AS
 SELECT a."ID" AS "ID_UNOR",
    a."UNOR_INDUK" AS "ID_SATKER",
    a."NAMA_UNOR",
    b."NAMA_UNOR" AS "NAMA_SATKER",
    c."NAMA_UNOR_ESELON_1",
    a."EXPIRED_DATE",
    c.id_eselon_1 AS "ID_ESELON_1"
   FROM ((unitkerja a
     JOIN unitkerja b ON (((a."UNOR_INDUK")::text = (b."ID")::text)))
     JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    unitkerja."ID" AS id_eselon_1,
                    unitkerja."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a_1."ID",
                    r_1.id_eselon_1,
                    r_1."NAMA_UNOR_ESELON_1"
                   FROM (unitkerja a_1
                     JOIN r r_1 ON (((a_1."DIATASAN_ID")::text = (r_1."ID")::text)))
                )
         SELECT r."ID",
            r.id_eselon_1,
            r."NAMA_UNOR_ESELON_1"
           FROM r) c ON (((a."ID")::text = (c."ID")::text)));

COMMENT ON VIEW vw_unor_satker IS 'Untuk Melihat Daftar Unit Kerja Berdasarkan Satkernya';

CREATE VIEW vw_rekap_input_diklat AS
 SELECT 'Diklat Fungsional'::text AS tipe,
    pegawai."NIP_BARU",
        CASE
            WHEN ((r_dik_fung."JUMLAH_JAM")::text = ''::text) THEN '0'::character varying
            ELSE r_dik_fung."JUMLAH_JAM"
        END AS "JUMLAH_JAM",
        CASE
            WHEN (split_part((r_dik_fung."TAHUN")::text, '-'::text, 1) ~ '^\d+$'::text) THEN (split_part((r_dik_fung."TAHUN")::text, '-'::text, 1))::integer
            ELSE 0
        END AS tahun,
    vw."NAMA_SATKER",
    pegawai."IS_DOSEN",
    vw."ID_SATKER"
   FROM ((rwt_diklat_fungsional r_dik_fung
     JOIN pegawai pegawai ON (((pegawai."NIP_BARU")::text = (r_dik_fung."NIP_BARU")::text)))
     JOIN vw_unor_satker vw ON (((pegawai."UNOR_ID")::text = (vw."ID_UNOR")::text)))
UNION ALL
 SELECT 'Diklat Struktural'::text AS tipe,
    pegawai."NIP_BARU",
    '20'::character varying AS "JUMLAH_JAM",
    (r_dik_struk."TAHUN")::integer AS tahun,
    vw."NAMA_SATKER",
    pegawai."IS_DOSEN",
    vw."ID_SATKER"
   FROM ((rwt_diklat_struktural r_dik_struk
     JOIN pegawai pegawai ON (((pegawai."NIP_BARU")::text = (r_dik_struk."PNS_NIP")::text)))
     JOIN vw_unor_satker vw ON (((pegawai."UNOR_ID")::text = (vw."ID_UNOR")::text)))
UNION ALL
 SELECT 'Riwayat Kursus'::text AS tipe,
    pegawai."NIP_BARU",
    (r_kurs."LAMA_KURSUS")::character varying AS "JUMLAH_JAM",
    (date_part('year'::text, r_kurs."TANGGAL_KURSUS"))::integer AS tahun,
    vw."NAMA_SATKER",
    pegawai."IS_DOSEN",
    vw."ID_SATKER"
   FROM ((rwt_kursus r_kurs
     JOIN pegawai pegawai ON (((pegawai."NIP_BARU")::bpchar = r_kurs."PNS_NIP")))
     JOIN vw_unor_satker vw ON (((pegawai."UNOR_ID")::text = (vw."ID_UNOR")::text)))
UNION ALL
 SELECT jd.jenis_diklat AS tipe,
    (r_dik.nip_baru)::character varying(18) AS "NIP_BARU",
    (r_dik.durasi_jam)::character varying AS "JUMLAH_JAM",
    r_dik.tahun_diklat AS tahun,
    vw."NAMA_SATKER",
    pegawai."IS_DOSEN",
    vw."ID_SATKER"
   FROM (((rwt_diklat r_dik
     JOIN pegawai pegawai ON (((pegawai."NIP_BARU")::text = (r_dik.nip_baru)::text)))
     JOIN vw_unor_satker vw ON (((pegawai."UNOR_ID")::text = (vw."ID_UNOR")::text)))
     JOIN jenis_diklat jd ON (((r_dik.jenis_diklat_id)::integer = jd.id)));

CREATE VIEW vw_rekap_pegawai_per_satker AS
 SELECT satker."ID_SATKER",
    satker."NAMA_SATKER",
    count(p."ID") AS total_pegawai
   FROM ((vw_unor_satker satker
     LEFT JOIN pegawai p ON (((satker."ID_UNOR")::text = (p."UNOR_ID")::text)))
     LEFT JOIN pns_aktif pa ON ((pa."ID" = p."ID")))
  WHERE (pa."ID" IS NOT NULL)
  GROUP BY satker."ID_SATKER", satker."NAMA_SATKER"
  ORDER BY satker."NAMA_SATKER" DESC;

COMMENT ON VIEW vw_rekap_pegawai_per_satker IS 'rekap pegawai aktif per satker';

CREATE VIEW vw_biro_sdm_award AS
 SELECT count(DISTINCT vw_diklat."NIP_BARU") AS total_pegawai_ngisi_diklat,
    vw_diklat."NAMA_SATKER",
    max(vw_rekap.total_pegawai) AS jumlah_pegawai_satker,
    round((((count(DISTINCT vw_diklat."NIP_BARU"))::numeric / max((vw_rekap.total_pegawai)::numeric)) * (100)::numeric)) AS percentage
   FROM (vw_rekap_pegawai_per_satker vw_rekap
     LEFT JOIN vw_rekap_input_diklat vw_diklat ON (((vw_rekap."ID_SATKER")::text = (vw_diklat."ID_SATKER")::text)))
  WHERE (vw_diklat.tahun = 2022)
  GROUP BY vw_diklat."NAMA_SATKER"
  ORDER BY vw_diklat."NAMA_SATKER";

CREATE VIEW vw_cuti_tahunan AS
 SELECT i."NIP_PNS",
    i."TAHUN",
    sum(i."JUMLAH") AS jumlah_hari
   FROM (izin i
     LEFT JOIN jenis_izin j ON ((i."KODE_IZIN" = j."ID")))
  WHERE ((i."KODE_IZIN" = 1) AND (i."STATUS_PENGAJUAN" = 3))
  GROUP BY i."KODE_IZIN", i."NIP_PNS", i."TAHUN";

CREATE VIEW vw_daftar_riwayat_jabatan AS
 SELECT row_number() OVER (PARTITION BY "PNS_ID" ORDER BY "TMT_JABATAN" DESC) AS _order,
    "ID_BKN",
    "PNS_ID",
    "PNS_NIP",
    "PNS_NAMA",
    "ID_UNOR",
    "UNOR",
    "ID_JENIS_JABATAN",
    "JENIS_JABATAN",
    "ID_JABATAN",
    "NAMA_JABATAN",
    "ID_ESELON",
    "ESELON",
    "TMT_JABATAN",
    "NOMOR_SK",
    "TANGGAL_SK",
    "ID_SATUAN_KERJA",
    "TMT_PELANTIKAN",
    "IS_ACTIVE",
    "ESELON1",
    "ESELON2",
    "ESELON3",
    "ESELON4",
    "ID",
    "CATATAN",
    "JENIS_SK",
    "LAST_UPDATED",
    "STATUS_SATKER",
    "STATUS_BIRO",
    "ID_JABATAN_BKN",
    "ID_UNOR_BKN",
    "JABATAN_TERAKHIR"
   FROM rwt_jabatan_empty rjab
  WHERE ("TMT_JABATAN" IS NOT NULL);

CREATE VIEW vw_drh AS
 SELECT pegawai."PNS_ID",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."GOL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."JABATAN_ID",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."TAHUN_LULUS",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."PHOTO",
    pegawai."LOKASI_KERJA",
    pegawai."TEMPAT_LAHIR",
    pegawai."PENDIDIKAN",
    pegawai."TK_PENDIDIKAN",
    pegawai."TEMPAT_LAHIR_NAMA",
    pegawai."JENIS_JABATAN_NAMA",
    pegawai."JABATAN_NAMA",
    pegawai."KPKN_NAMA",
    pegawai."INSTANSI_INDUK_NAMA",
    pegawai."INSTANSI_KERJA_NAMA",
    pegawai."SATUAN_KERJA_INDUK_NAMA",
    pegawai."SATUAN_KERJA_NAMA",
    pegawai."BUP",
    golongan."NAMA" AS "GOL_TEXT",
    golongan."NAMA_PANGKAT" AS "PANGKAT_TEXT",
    agama."NAMA" AS "AGAMA_TEXT",
    jenis_kawin."NAMA" AS "KAWIN_TEXT"
   FROM (((pegawai pegawai
     LEFT JOIN golongan ON ((golongan."ID" = pegawai."GOL_ID")))
     LEFT JOIN agama ON ((agama."ID" = pegawai."AGAMA_ID")))
     LEFT JOIN jenis_kawin ON (((jenis_kawin."ID")::text = (pegawai."JENIS_KAWIN_ID")::text)));

CREATE VIEW vw_ds_korektor AS
 SELECT d.id_file,
    k.id_pegawai_korektor,
    k.is_corrected,
    d.nomor_sk
   FROM (tbl_file_ds d
     JOIN tbl_file_ds_corrector k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND (k.is_corrected <> (1)::smallint));

CREATE VIEW vw_ds_antrian_korektor AS
 SELECT d.kategori,
    k.id_pegawai_korektor,
    count(*) AS jumlah
   FROM (tbl_file_ds d
     JOIN vw_ds_korektor k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text))
  GROUP BY d.kategori, k.id_pegawai_korektor;

CREATE VIEW vw_ds_jml_korektor_new_1 AS
 SELECT d.id_file,
    count(k.id) AS jumlah_korektor
   FROM (tbl_file_ds d
     JOIN tbl_file_ds_corrector k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1))
  GROUP BY d.id_file;

CREATE VIEW vw_ds_antrian_ttd AS
 SELECT d.kategori,
    d.id_pegawai_ttd,
    d.is_signed,
    count(*) AS jumlah
   FROM (tbl_file_ds d
     JOIN vw_ds_jml_korektor_new_1 k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND ((d.kategori)::text <> '< Pilih >'::text) AND (k.jumlah_korektor > 0))
  GROUP BY d.kategori, d.is_signed, d.id_pegawai_ttd;

CREATE VIEW vw_ds_antrian_ttd_copy1 AS
 SELECT d.kategori,
    d.id_pegawai_ttd,
    d.is_signed,
    count(*) AS jumlah
   FROM (tbl_file_ds d
     JOIN vw_ds_jml_korektor_new_1 k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND ((d.kategori)::text <> '< Pilih >'::text) AND (k.jumlah_korektor > 0))
  GROUP BY d.kategori, d.is_signed, d.id_pegawai_ttd;

CREATE VIEW vw_ds_jml_korektor AS
 SELECT d.id_file,
    count(k.id) AS jumlah_korektor
   FROM (tbl_file_ds d
     JOIN tbl_file_ds_corrector k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND (d.is_signed <> (3)::smallint))
  GROUP BY d.id_file;

CREATE VIEW vw_ds_jml_korektor_new AS
 SELECT d.id_file,
    count(k.id) AS jumlah_korektor
   FROM (tbl_file_ds d
     JOIN tbl_file_ds_corrector k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE (d.ds_ok = 1)
  GROUP BY d.id_file;

CREATE VIEW vw_ds_jumlah_pernip AS
 SELECT d.nip_sk,
    d.is_signed,
    count(*) AS jumlah
   FROM (tbl_file_ds d
     JOIN vw_ds_jml_korektor_new k ON (((d.id_file)::text = (k.id_file)::text)))
  WHERE ((d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND (k.jumlah_korektor > 0))
  GROUP BY d.is_signed, d.nip_sk;

CREATE VIEW vw_ds_pejabat_ttd_dan_korektor AS
 SELECT "PNS_ID",
    "NAMA"
   FROM pegawai d
  WHERE ((("PNS_ID")::text IN ( SELECT tbl_file_ds.id_pegawai_ttd
           FROM tbl_file_ds
          WHERE ((tbl_file_ds.ds_ok = '1'::smallint) AND (tbl_file_ds.is_signed <> '1'::smallint) AND (tbl_file_ds.is_signed <> '3'::smallint)))) OR (("PNS_ID")::text IN ( SELECT tbl_file_ds_corrector.id_pegawai_korektor
           FROM tbl_file_ds_corrector
          WHERE ((tbl_file_ds_corrector.is_corrected = '2'::smallint) AND ((tbl_file_ds_corrector.id_file)::text IN ( SELECT tbl_file_ds.id_file
                   FROM tbl_file_ds
                  WHERE ((tbl_file_ds.ds_ok = '1'::smallint) AND (tbl_file_ds.is_signed <> '1'::smallint) AND (tbl_file_ds.is_signed <> '3'::smallint))))))));

CREATE VIEW vw_ds_siap_koreksi AS
 SELECT k.id_pegawai_korektor,
    p."NAMA" AS nama_korektor,
    count(*) AS jumlah
   FROM ((tbl_file_ds d
     JOIN tbl_file_ds_corrector k ON (((d.id_file)::text = (k.id_file)::text)))
     JOIN pegawai p ON (((p."PNS_ID")::text = (k.id_pegawai_korektor)::text)))
  WHERE ((d.is_signed = (0)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND (k.is_corrected = 2))
  GROUP BY k.id_pegawai_korektor, p."NAMA";

CREATE VIEW vw_ds_siap_ttd AS
 SELECT d.id_pegawai_ttd,
    p."NAMA" AS nama_penandatangan,
    count(*) AS jumlah
   FROM (tbl_file_ds d
     JOIN pegawai p ON (((p."PNS_ID")::text = (d.id_pegawai_ttd)::text)))
  WHERE ((d.is_signed = (0)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text) AND (d.is_corrected = 1))
  GROUP BY d.id_pegawai_ttd, p."NAMA";

CREATE VIEW vw_ds_resume_ttd AS
 SELECT d.id_pegawai_ttd,
    p."NAMA" AS nama_penandatangan,
    COALESCE(( SELECT (sk.jumlah)::text AS jumlah
           FROM vw_ds_siap_koreksi sk
          WHERE ((sk.id_pegawai_korektor)::text = (d.id_pegawai_ttd)::text)), '-'::text) AS jml_siap_koreksi,
    COALESCE(( SELECT (s.jumlah)::text AS jumlah
           FROM vw_ds_siap_ttd s
          WHERE ((s.id_pegawai_ttd)::text = (d.id_pegawai_ttd)::text)), '-'::text) AS jml_siap_ttd,
    (COALESCE(( SELECT sk.jumlah
           FROM vw_ds_siap_koreksi sk
          WHERE ((sk.id_pegawai_korektor)::text = (d.id_pegawai_ttd)::text)), (0)::bigint) + COALESCE(( SELECT s.jumlah
           FROM vw_ds_siap_ttd s
          WHERE ((s.id_pegawai_ttd)::text = (d.id_pegawai_ttd)::text)), (0)::bigint)) AS jumlah
   FROM (tbl_file_ds d
     JOIN pegawai p ON (((p."PNS_ID")::text = (d.id_pegawai_ttd)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text))
  GROUP BY d.id_pegawai_ttd, p."NAMA";

CREATE VIEW vw_ds_resume_ttd_copy1 AS
 SELECT d.id_pegawai_ttd,
    p."NAMA" AS nama_penandatangan,
    count(*) AS jumlah,
    ( SELECT s.jumlah
           FROM vw_ds_siap_ttd s
          WHERE ((s.id_pegawai_ttd)::text = (d.id_pegawai_ttd)::text)) AS jml_siap_ttd,
    ( SELECT sk.jumlah
           FROM vw_ds_siap_koreksi sk
          WHERE ((sk.id_pegawai_korektor)::text = (d.id_pegawai_ttd)::text)) AS jml_siap_koreksi
   FROM ((tbl_file_ds d
     JOIN vw_ds_jml_korektor k ON (((d.id_file)::text = (k.id_file)::text)))
     JOIN pegawai p ON (((p."PNS_ID")::text = (d.id_pegawai_ttd)::text)))
  WHERE ((d.is_signed <> (1)::smallint) AND (d.ds_ok = 1) AND ((d.kategori)::text <> '< Semua >'::text))
  GROUP BY d.id_pegawai_ttd, p."NAMA";

CREATE VIEW vw_ds_resume_ttd_copy3 AS
 SELECT "PNS_ID" AS id_pegawai_ttd,
    "NAMA" AS nama_penandatangan,
    COALESCE(( SELECT (sk.jumlah)::text AS jumlah
           FROM vw_ds_siap_koreksi sk
          WHERE ((sk.id_pegawai_korektor)::text = (d."PNS_ID")::text)), '-'::text) AS jml_siap_koreksi,
    COALESCE(( SELECT (s.jumlah)::text AS jumlah
           FROM vw_ds_siap_ttd s
          WHERE ((s.id_pegawai_ttd)::text = (d."PNS_ID")::text)), '-'::text) AS jml_siap_ttd,
    '-'::text AS jumlah
   FROM vw_ds_pejabat_ttd_dan_korektor d;

CREATE VIEW vw_duk AS
 SELECT vw."NAMA_UNOR",
    pegawai."JENIS_JABATAN_ID",
    pegawai."JABATAN_ID",
    pegawai."JABATAN_NAMA",
    pegawai."NIP_LAMA",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    vw."ESELON_ID" AS vw_eselon_id,
    pegawai."GOL_ID",
    (((golongan."NAMA_PANGKAT")::text || ' '::text) || (golongan."NAMA")::text) AS golongan_text,
    'jabatanku'::text AS jabatan_text,
    pegawai."PNS_ID",
    (((date_part('year'::text, (now())::date) - date_part('year'::text, pegawai."TGL_LAHIR")) * (12)::double precision) + (date_part('month'::text, (now())::date) - date_part('month'::text, pegawai."TGL_LAHIR"))) AS bulan_usia,
    '#'::text AS separator,
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."KEDUDUKAN_HUKUM_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."KARTU_PEGAWAI",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."PENDIDIKAN",
    pegawai."TAHUN_LULUS",
    pegawai."KPKN_ID",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."ID",
    pegawai."PHOTO",
    pegawai."TMT_PENSIUN",
    pegawai."BUP",
    pegawai."EMAIL_DIKBUD",
    vw."ESELON_1",
    vw."ESELON_2",
    vw."ESELON_3",
    vw."ESELON_4",
    tkpendidikan."NAMA" AS "TINGKAT_PENDIDIKAN_NAMA",
    jabatan."NAMA_JABATAN",
    jabatan."KELAS",
    jabatan."PENSIUN",
    kedudukan_hukum."NAMA" AS "KEDUDUKAN_HUKUM_NAMA"
   FROM (((((((pns_aktif pa
     LEFT JOIN pegawai pegawai ON ((pa."ID" = pegawai."ID")))
     LEFT JOIN golongan ON (((pegawai."GOL_ID")::text = (golongan."ID")::text)))
     LEFT JOIN vw_unit_list vw ON (((vw."ID")::text = (pegawai."UNOR_ID")::text)))
     LEFT JOIN pendidikan ON (((pendidikan."ID")::text = (pegawai."PENDIDIKAN_ID")::text)))
     LEFT JOIN tkpendidikan ON (((tkpendidikan."ID")::text = (pendidikan."TINGKAT_PENDIDIKAN_ID")::text)))
     LEFT JOIN jabatan ON (((jabatan."KODE_JABATAN")::text = (pegawai."JABATAN_INSTANSI_ID")::text)))
     LEFT JOIN kedudukan_hukum ON (((kedudukan_hukum."ID")::text = (pegawai."KEDUDUKAN_HUKUM_ID")::text)))
  ORDER BY pegawai."JENIS_JABATAN_ID", vw."ESELON_ID", vw."ESELON_1", vw."ESELON_2", vw."ESELON_3", vw."ESELON_4", pegawai."JABATAN_ID", vw."NAMA_UNOR_FULL", pegawai."GOL_ID" DESC, pegawai."TMT_GOLONGAN", pegawai."TMT_JABATAN", pegawai."TMT_CPNS", pegawai."TGL_LAHIR";

CREATE VIEW vw_duk_list AS
 SELECT vw."NAMA_UNOR_FULL",
    pejabat."ESELON_ID",
    vw."ESELON_ID" AS vw_eselon_id,
    pegawai."GOL_ID",
    (((golongan."NAMA_PANGKAT")::text || ' '::text) || (golongan."NAMA")::text) AS golongan_text,
    'jabatanku'::text AS jabatan_text,
    pegawai."PNS_ID",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    (((date_part('year'::text, (now())::date) - date_part('year'::text, pegawai."TGL_LAHIR")) * (12)::double precision) + (date_part('month'::text, (now())::date) - date_part('month'::text, pegawai."TGL_LAHIR"))) AS bulan_usia,
    '#'::text AS separator,
    pegawai."NIP_LAMA",
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."KEDUDUKAN_HUKUM_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."KARTU_PEGAWAI",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."JENIS_JABATAN_IDx" AS "JENIS_JABATAN_ID",
    pegawai."JABATAN_ID",
    pegawai."JABATAN_NAMA",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."PENDIDIKAN",
    pegawai."TAHUN_LULUS",
    pegawai."KPKN_ID",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."ID",
    pegawai."PHOTO",
    pegawai."TMT_PENSIUN",
    pegawai."BUP",
    vw."NAMA_UNOR",
    vw."ESELON_1",
    vw."ESELON_2",
    vw."ESELON_3",
    vw."ESELON_4"
   FROM ((((pns_aktif_old pa
     LEFT JOIN pegawai pegawai ON ((pa."ID" = pegawai."ID")))
     LEFT JOIN golongan ON (((pegawai."GOL_ID")::text = (golongan."ID")::text)))
     LEFT JOIN vw_unit_list vw ON (((vw."ID")::text = (pegawai."UNOR_ID")::text)))
     LEFT JOIN unitkerja pejabat ON (((pejabat."PEMIMPIN_PNS_ID")::text = (pegawai."PNS_ID")::text)))
  ORDER BY vw."NAMA_UNOR_FULL", pejabat."ESELON_ID", pegawai."GOL_ID" DESC, pegawai."TMT_GOLONGAN", pegawai."TMT_JABATAN", pegawai."TMT_CPNS", pegawai."TGL_LAHIR";

CREATE VIEW vw_kgb AS
 SELECT get_masa_kerja_arr("TMT_CPNS", get_kgb_yad) AS get_masa_kerja_arr,
    get_kgb_yad,
    "ID",
    "PNS_ID",
    "NIP_LAMA",
    "NIP_BARU",
    "NAMA",
    "GELAR_DEPAN",
    "GELAR_BELAKANG",
    "TEMPAT_LAHIR_ID",
    "TGL_LAHIR",
    "JENIS_KELAMIN",
    "AGAMA_ID",
    "JENIS_KAWIN_ID",
    "NIK",
    "NOMOR_DARURAT",
    "NOMOR_HP",
    "EMAIL",
    "ALAMAT",
    "NPWP",
    "BPJS",
    "JENIS_PEGAWAI_ID",
    "KEDUDUKAN_HUKUM_ID",
    "STATUS_CPNS_PNS",
    "KARTU_PEGAWAI",
    "NOMOR_SK_CPNS",
    "TGL_SK_CPNS",
    "TMT_CPNS",
    "TMT_PNS",
    "GOL_AWAL_ID",
    "GOL_ID",
    "TMT_GOLONGAN",
    "MK_TAHUN",
    "MK_BULAN",
    "JENIS_JABATAN_IDx",
    "JABATAN_ID",
    "TMT_JABATAN",
    "PENDIDIKAN_ID",
    "TAHUN_LULUS",
    "KPKN_ID",
    "LOKASI_KERJA_ID",
    "UNOR_ID",
    "UNOR_INDUK_ID",
    "INSTANSI_INDUK_ID",
    "INSTANSI_KERJA_ID",
    "SATUAN_KERJA_INDUK_ID",
    "SATUAN_KERJA_KERJA_ID",
    "GOLONGAN_DARAH",
    "PHOTO",
    "TMT_PENSIUN",
    "LOKASI_KERJA",
    "JML_ISTRI",
    "JML_ANAK",
    "NO_SURAT_DOKTER",
    "TGL_SURAT_DOKTER",
    "NO_BEBAS_NARKOBA",
    "TGL_BEBAS_NARKOBA",
    "NO_CATATAN_POLISI",
    "TGL_CATATAN_POLISI",
    "AKTE_KELAHIRAN",
    "STATUS_HIDUP",
    "AKTE_MENINGGAL",
    "TGL_MENINGGAL",
    "NO_ASKES",
    "NO_TASPEN",
    "TGL_NPWP",
    "TEMPAT_LAHIR",
    "PENDIDIKAN",
    "TK_PENDIDIKAN",
    "TEMPAT_LAHIR_NAMA",
    "JENIS_JABATAN_NAMA",
    "JABATAN_NAMA",
    "KPKN_NAMA",
    "INSTANSI_INDUK_NAMA",
    "INSTANSI_KERJA_NAMA",
    "SATUAN_KERJA_INDUK_NAMA",
    "SATUAN_KERJA_NAMA",
    "JABATAN_INSTANSI_ID",
    "BUP",
    "JABATAN_INSTANSI_NAMA",
    "JENIS_JABATAN_ID",
    terminated_date,
    status_pegawai
   FROM ( SELECT get_kgb_yad(p."TMT_CPNS") AS get_kgb_yad,
            p."ID",
            p."ID" AS pegawai_id,
            p."PNS_ID",
            p."NIP_LAMA",
            p."NIP_BARU",
            p."NAMA",
            p."GELAR_DEPAN",
            p."GELAR_BELAKANG",
            p."TEMPAT_LAHIR_ID",
            p."TGL_LAHIR",
            p."JENIS_KELAMIN",
            p."AGAMA_ID",
            p."JENIS_KAWIN_ID",
            p."NIK",
            p."NOMOR_DARURAT",
            p."NOMOR_HP",
            p."EMAIL",
            p."ALAMAT",
            p."NPWP",
            p."BPJS",
            p."JENIS_PEGAWAI_ID",
            p."KEDUDUKAN_HUKUM_ID",
            p."STATUS_CPNS_PNS",
            p."KARTU_PEGAWAI",
            p."NOMOR_SK_CPNS",
            p."TGL_SK_CPNS",
            p."TMT_CPNS",
            p."TMT_PNS",
            p."GOL_AWAL_ID",
            p."GOL_ID",
            p."TMT_GOLONGAN",
            p."MK_TAHUN",
            p."MK_BULAN",
            p."JENIS_JABATAN_IDx",
            p."JABATAN_ID",
            p."TMT_JABATAN",
            p."PENDIDIKAN_ID",
            p."TAHUN_LULUS",
            p."KPKN_ID",
            p."LOKASI_KERJA_ID",
            p."UNOR_ID",
            p."UNOR_INDUK_ID",
            p."INSTANSI_INDUK_ID",
            p."INSTANSI_KERJA_ID",
            p."SATUAN_KERJA_INDUK_ID",
            p."SATUAN_KERJA_KERJA_ID",
            p."GOLONGAN_DARAH",
            p."PHOTO",
            p."TMT_PENSIUN",
            p."LOKASI_KERJA",
            p."JML_ISTRI",
            p."JML_ANAK",
            p."NO_SURAT_DOKTER",
            p."TGL_SURAT_DOKTER",
            p."NO_BEBAS_NARKOBA",
            p."TGL_BEBAS_NARKOBA",
            p."NO_CATATAN_POLISI",
            p."TGL_CATATAN_POLISI",
            p."AKTE_KELAHIRAN",
            p."STATUS_HIDUP",
            p."AKTE_MENINGGAL",
            p."TGL_MENINGGAL",
            p."NO_ASKES",
            p."NO_TASPEN",
            p."TGL_NPWP",
            p."TEMPAT_LAHIR",
            p."PENDIDIKAN",
            p."TK_PENDIDIKAN",
            p."TEMPAT_LAHIR_NAMA",
            p."JENIS_JABATAN_NAMA",
            p."JABATAN_NAMA",
            p."KPKN_NAMA",
            p."INSTANSI_INDUK_NAMA",
            p."INSTANSI_KERJA_NAMA",
            p."SATUAN_KERJA_INDUK_NAMA",
            p."SATUAN_KERJA_NAMA",
            p."JABATAN_INSTANSI_ID",
            p."BUP",
            p."JABATAN_INSTANSI_NAMA",
            p."JENIS_JABATAN_ID",
            p.terminated_date,
            p.status_pegawai
           FROM (pns_aktif_old pa
             JOIN pegawai p ON ((pa."ID" = p."ID")))) temp;

CREATE VIEW vw_list_eselon1 AS
 SELECT DISTINCT "NO",
    "KODE_INTERNAL",
    "ID",
    "NAMA_UNOR",
    "ESELON_ID",
    "CEPAT_KODE",
    "NAMA_JABATAN",
    "NAMA_PEJABAT",
    "DIATASAN_ID",
    "INSTANSI_ID",
    "PEMIMPIN_NON_PNS_ID",
    "PEMIMPIN_PNS_ID",
    "JENIS_UNOR_ID",
    "UNOR_INDUK",
    "JUMLAH_IDEAL_STAFF",
    "ORDER",
    deleted,
    "IS_SATKER",
    "ESELON_1",
    "ESELON_2",
    "ESELON_3",
    "ESELON_4",
    "EXPIRED_DATE",
    "KETERANGAN",
    "ABBREVIATION"
   FROM unitkerja es1
  WHERE (("ID" IS NOT NULL) AND (("DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text) AND ("EXPIRED_DATE" IS NULL))
  ORDER BY "NAMA_UNOR";

CREATE VIEW vw_list_eselon2 AS
 SELECT DISTINCT es2."NO",
    es2."KODE_INTERNAL",
    es2."ID",
    es2."NAMA_UNOR",
    es2."ESELON_ID",
    es2."CEPAT_KODE",
    es2."NAMA_JABATAN",
    es2."NAMA_PEJABAT",
    es2."DIATASAN_ID",
    es2."INSTANSI_ID",
    es2."PEMIMPIN_NON_PNS_ID",
    es2."PEMIMPIN_PNS_ID",
    es2."JENIS_UNOR_ID",
    es2."UNOR_INDUK",
    es2."JUMLAH_IDEAL_STAFF",
    es2."ORDER",
    es2.deleted,
    es2."IS_SATKER",
    es2."ESELON_1",
    es2."ESELON_2",
    es2."ESELON_3",
    es2."ESELON_4",
    es2."EXPIRED_DATE",
    es2."KETERANGAN",
    es2."ABBREVIATION"
   FROM (unitkerja uk
     LEFT JOIN unitkerja es2 ON (((uk."ESELON_2")::text = (es2."ID")::text)))
  WHERE (es2."ID" IS NOT NULL)
  ORDER BY es2."NAMA_UNOR";

CREATE VIEW vw_pegawai_berpotensi_jpt AS
 SELECT apbt.id,
    apbt.nip,
    apbt.usia,
    apbt.status_kepegawaian,
    apbt.golongan,
    apbt.jenis_jabatan,
    apbt.jabatan,
    apbt.tmt,
    apbt.lama_jabatan_terakhir,
    apbt.eselon,
    apbt.satker,
    apbt.unit_organisasi_induk,
    apbt.kedudukan,
    apbt.tipe,
    apbt.pendidikan,
    apbt.jabatan_madya_lain,
    apbt.tmt_jabatan_madya_lain,
    apbt.jabatan_struktural_lain,
    apbt.tmt_jabatan_struktural_lain,
    apbt.lama_menjabat_akumulasi,
    apbt.rekam_jejak,
    apbt.skp,
    tp."NAMA" AS tingkat_pendidikan,
    pd."NAMA" AS nama_pendidikan,
    aha.tanggal_asesmen,
    aha.jpm
   FROM ((((asesmen_pegawai_berpotensi_jpt apbt
     LEFT JOIN asesmen_hasil_asesmen aha ON ((btrim((apbt.nip)::text) = btrim((aha.nip)::text))))
     LEFT JOIN pegawai p ON (((p."NIP_BARU")::text = (apbt.nip)::text)))
     LEFT JOIN pendidikan pd ON (((p."PENDIDIKAN_ID")::text = (pd."ID")::text)))
     LEFT JOIN tkpendidikan tp ON (((tp."ID")::text = (pd."TINGKAT_PENDIDIKAN_ID")::text)));

CREATE VIEW vw_pegawai_bpk AS
 SELECT pegawai."NAMA" AS "NAMA_PEGAWAI",
    (''''::text || (pegawai."NIP_BARU")::text) AS "NIP_BARU",
    (''''::text || (pegawai."NIK")::text) AS "NIK",
    lokasi."NAMA" AS "TEMPAT_LAHIR",
    pegawai."TGL_LAHIR" AS "TANGGAL_LAHIR",
    agama."NAMA" AS "AGAMA",
    pegawai."JENIS_KELAMIN",
    tkpendidikan."NAMA" AS "TINGKAT_PENDIDIKAN",
    pendidikan."NAMA" AS "NAMA_PENDIDIKAN",
    (((pegawai."MK_TAHUN")::text || '/'::text) || (pegawai."MK_BULAN")::text) AS "MASA_KERJA",
    golongan."NAMA" AS "PANGKAT_GOLONGAN_AKTIF",
    golongan_awal."NAMA" AS "GOLONGAN_AWAL",
    'PNS PUSAT'::text AS "JENIS_PEGAWAI",
    pegawai."TMT_GOLONGAN",
    pegawai."TMT_CPNS",
    pegawai."STATUS_CPNS_PNS",
    jenis_kawin."NAMA" AS "JENIS_KAWIN",
    pegawai."NPWP",
    pegawai."EMAIL",
    pegawai."EMAIL_DIKBUD",
    pegawai."NOMOR_HP",
    jabatan."NAMA_JABATAN",
    jabatan."KATEGORI_JABATAN",
    vw."NAMA_UNOR",
    vw."NAMA_UNOR_FULL"
   FROM ((((((((((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
     LEFT JOIN jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
     LEFT JOIN lokasi ON (((pegawai."TEMPAT_LAHIR_ID")::text = (lokasi."ID")::text)))
     LEFT JOIN agama ON ((pegawai."AGAMA_ID" = agama."ID")))
     LEFT JOIN tkpendidikan ON ((pegawai."TK_PENDIDIKAN" = (tkpendidikan."ID")::bpchar)))
     LEFT JOIN pendidikan ON (((pegawai."PENDIDIKAN_ID")::text = (pendidikan."ID")::text)))
     LEFT JOIN golongan ON ((pegawai."GOL_ID" = golongan."ID")))
     LEFT JOIN golongan golongan_awal ON (((pegawai."GOL_AWAL_ID")::text = (golongan_awal."ID")::text)))
     LEFT JOIN jenis_kawin ON (((jenis_kawin."ID")::text = (pegawai."JENIS_KAWIN_ID")::text)))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('88'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text, ('100'::character varying)::text])) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)));

CREATE VIEW vw_pegawai_for_bpjstk AS
 SELECT mvp."NIK",
    mvp."NAMA" AS nama_pegawai,
    (((l."NAMA")::text || ' / '::text) || mvp."TGL_LAHIR") AS tempat_tanggal_lahir,
    mvp."NOMOR_HP" AS nomor_telepon,
    mvp."EMAIL" AS email,
    mvp."NAMA_UNOR_FULL" AS unit_kerja,
    mvp."NAMA_JABATAN_REAL" AS jabatan,
    mvp."KATEGORI_JABATAN_REAL" AS jenis_jabatan,
    mvp."NAMA_GOLONGAN" AS golongan,
        CASE
            WHEN ((mvp."JENIS_PEGAWAI_ID" = '71'::text) OR (mvp."JENIS_PEGAWAI_ID" = '72'::text) OR (mvp."JENIS_PEGAWAI_ID" = '73'::text)) THEN 'PPPK'::text
            ELSE 'PNS'::text
        END AS status_kepegawaian,
    mvp."UNOR_ID",
    mvp."UNOR_INDUK_ID"
   FROM (mv_pegawai mvp
     LEFT JOIN lokasi l ON ((mvp."TEMPAT_LAHIR_ID" = (l."ID")::text)));

CREATE VIEW vw_pegawai_simple AS
 SELECT pegawai."ID",
    pegawai."NIP_BARU",
    btrim((pegawai."NAMA")::text) AS "NAMA",
    btrim((pegawai."UNOR_INDUK_ID")::text) AS "UNOR_INDUK_ID",
    btrim((pegawai."UNOR_ID")::text) AS "UNOR_ID",
    vw."ESELON_1" AS "VW_ESELON_1",
    vw."ESELON_2" AS "VW_ESELON_2",
    vw."ESELON_3" AS "VW_ESELON_3",
    vw."UNOR_INDUK" AS "VW_UNOR_INDUK"
   FROM ((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '99'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '66'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '52'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '20'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '04'::text) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)));

CREATE VIEW vw_pegawai_tanpa_akun AS
 SELECT pegawai."NIP_BARU"
   FROM (((pegawai pegawai
     LEFT JOIN vw_unit_list vw ON (((pegawai."UNOR_ID")::text = (vw."ID")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
     LEFT JOIN jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
  WHERE ((pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('78'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text])) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)) AND (NOT ((pegawai."NIP_BARU")::text IN ( SELECT users.username
           FROM users))));

CREATE VIEW vw_pejabat_cuti AS
 SELECT "NIP_ATASAN",
    count(*) AS jumlah
   FROM line_approval_izin
  GROUP BY "NIP_ATASAN";

CREATE VIEW vw_rwt_asesmen_terakhir AS
 SELECT concat('|', pegawai."NIP_BARU") AS "NIP_BARU",
    pegawai."NAMA",
    pegawai."JENIS_KELAMIN",
    golongan."NAMA" AS golongan,
    pegawai."UNOR_ID" AS "UNOR_PEGAWAI",
    pegawai."SATUAN_KERJA_INDUK_ID" AS "SATKER_PEGAWAI",
    pegawai."JABATAN_NAMA",
    pegawai."JENIS_JABATAN_NAMA",
    vw_unor_satker."NAMA_UNOR",
    vw_unor_satker."NAMA_SATKER",
    vw_unor_satker."NAMA_UNOR_ESELON_1",
    rwt_assesmen."ID",
    rwt_assesmen."PNS_NIP",
    max((rwt_assesmen."TAHUN")::integer) AS last_asesmen
   FROM ((((pegawai pegawai
     LEFT JOIN golongan ON ((pegawai."GOL_ID" = golongan."ID")))
     LEFT JOIN rwt_assesmen ON (((pegawai."NIP_BARU")::bpchar = rwt_assesmen."PNS_NIP")))
     LEFT JOIN vw_unor_satker ON (((pegawai."UNOR_ID")::text = (vw_unor_satker."ID_UNOR")::text)))
     LEFT JOIN pns_aktif pa ON ((pegawai."ID" = pa."ID")))
  WHERE ((pegawai.status_pegawai = 1) AND ((pegawai.terminated_date IS NULL) OR ((pegawai.terminated_date IS NOT NULL) AND (pegawai.terminated_date > ('now'::text)::date))) AND (pa."ID" IS NOT NULL) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '99'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '66'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '52'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '20'::text) AND ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> '04'::text) AND ((pegawai.status_pegawai <> 3) OR (pegawai.status_pegawai IS NULL)))
  GROUP BY pegawai."NIP_BARU", pegawai."NAMA", pegawai."JENIS_KELAMIN", golongan."NAMA", pegawai."UNOR_ID", pegawai."SATUAN_KERJA_INDUK_ID", pegawai."JABATAN_NAMA", pegawai."JENIS_JABATAN_NAMA", vw_unor_satker."NAMA_UNOR", vw_unor_satker."NAMA_SATKER", vw_unor_satker."NAMA_UNOR_ESELON_1", rwt_assesmen."ID", rwt_assesmen."PNS_ID", rwt_assesmen."PNS_NIP";

CREATE VIEW vw_rwt_assesmen_pegawai AS
 SELECT pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."UNOR_ID" AS "UNOR_PEGAWAI",
    pegawai."UNOR_INDUK_ID" AS "UNOR_INDUK_PEGAWAI",
    pegawai."SATUAN_KERJA_INDUK_ID" AS "SATKER_PEGAWAI",
    rwt_assesmen."ID",
    rwt_assesmen."PNS_ID",
    rwt_assesmen."PNS_NIP",
    rwt_assesmen."TAHUN",
    rwt_assesmen."FILE_UPLOAD",
    rwt_assesmen."NILAI",
    rwt_assesmen."NILAI_KINERJA",
    rwt_assesmen."TAHUN_PENILAIAN_ID",
    rwt_assesmen."TAHUN_PENILAIAN_TITLE",
    rwt_assesmen."FULLNAME",
    rwt_assesmen."POSISI_ID",
    rwt_assesmen."UNIT_ORG_ID",
    rwt_assesmen."NAMA_UNOR",
    rwt_assesmen."SARANPENGEMBANGAN",
    rwt_assesmen."FILE_UPLOAD_FB_POTENSI",
    rwt_assesmen."FILE_UPLOAD_LENGKAP_PT",
    rwt_assesmen."FILE_UPLOAD_FB_PT",
    rwt_assesmen."FILE_UPLOAD_EXISTS",
    rwt_assesmen."SATKER_ID",
    vw_unor_satker."NAMA_UNOR" AS "NAMA_UNOR_PEGAWAI"
   FROM ((pegawai pegawai
     LEFT JOIN rwt_assesmen ON (((pegawai."NIP_BARU")::bpchar = rwt_assesmen."PNS_NIP")))
     LEFT JOIN vw_unor_satker ON (((pegawai."UNOR_ID")::text = (vw_unor_satker."ID_UNOR")::text)))
  WHERE ((rwt_assesmen."ID" IS NOT NULL) AND (pegawai.status_pegawai = 1) AND ((pegawai.terminated_date IS NULL) OR ((pegawai.terminated_date IS NOT NULL) AND (pegawai.terminated_date > ('now'::text)::date))));

CREATE VIEW vw_rwt_diklat AS
 SELECT rwt_diklat.id,
    jenis_diklat_siasn.jenis_diklat,
    jenis_diklat_siasn.id AS jenis_diklat_id,
    rwt_diklat.institusi_penyelenggara,
    rwt_diklat.nomor_sertifikat,
    rwt_diklat.tanggal_mulai,
    rwt_diklat.tanggal_selesai,
    rwt_diklat.tahun_diklat,
    rwt_diklat.durasi_jam,
    rwt_diklat.nip_baru,
    rwt_diklat.createddate,
    rwt_diklat.nama_diklat,
    jenis_rumpun_diklat_siasn.nama AS rumpun_diklat,
    jenis_rumpun_diklat_siasn.id AS rumpun_diklat_id,
    rwt_diklat.sudah_kirim_siasn,
    pegawai."PNS_ID" AS pns_orang_id,
    rwt_diklat.siasn_id,
    rwt_diklat.diklat_struktural_id
   FROM (((rwt_diklat
     LEFT JOIN jenis_diklat_siasn ON (((rwt_diklat.jenis_diklat_id)::text = (jenis_diklat_siasn.id)::text)))
     LEFT JOIN jenis_rumpun_diklat_siasn ON (((jenis_rumpun_diklat_siasn.id)::text = (rwt_diklat.rumpun_diklat_id)::text)))
     LEFT JOIN pegawai ON (((pegawai."NIP_BARU")::text = (rwt_diklat.nip_baru)::text)));

CREATE VIEW vw_skp AS
 SELECT s."ID",
    s."PNS_NIP",
    s."PNS_NAMA",
    s."ATASAN_LANGSUNG_PNS_NAMA",
    s."ATASAN_LANGSUNG_PNS_NIP",
    s."NILAI_SKP",
    s."NILAI_PROSENTASE_SKP",
    s."NILAI_SKP_AKHIR",
    s."PERILAKU_KOMITMEN",
    s."PERILAKU_INTEGRITAS",
    s."PERILAKU_DISIPLIN",
    s."PERILAKU_KERJASAMA",
    s."PERILAKU_ORIENTASI_PELAYANAN",
    s."PERILAKU_KEPEMIMPINAN",
    s."NILAI_PERILAKU",
    s."NILAI_PROSENTASE_PERILAKU",
    s."NILAI_PERILAKU_AKHIR",
    s."NILAI_PPK",
    s."TAHUN",
    s."JABATAN_TIPE",
    s."PNS_ID",
    s."ATASAN_LANGSUNG_PNS_ID",
    s."ATASAN_ATASAN_LANGSUNG_PNS_ID",
    s."ATASAN_ATASAN_LANGSUNG_PNS_NAMA",
    s."ATASAN_ATASAN_LANGSUNG_PNS_NIP",
    s."JABATAN_TIPE_TEXT",
    s."ATASAN_LANGSUNG_PNS_JABATAN",
    s."ATASAN_ATASAN_LANGSUNG_PNS_JABATAN",
    s."JABATAN_NAMA",
    s."BKN_ID",
    s."UNOR_PENILAI",
    s."UNOR_ATASAN_PENILAI",
    s."ATASAN_PENILAI_PNS",
    s."PENILAI_PNS",
    s."GOL_PENILAI",
    s."GOL_ATASAN_PENILAI",
    s."TMT_GOL_PENILAI",
    s."TMT_GOL_ATASAN_PENILAI",
    s."PERATURAN",
    s.created_date,
    s.updated_date,
    s."PERILAKU_INISIATIF_KERJA",
    a."NIP_BARU" AS nip_atasan,
    a."NAMA" AS nama_atasan,
    a."PNS_ID" AS pns_id_atasan,
    aa."NIP_BARU" AS a_nip_atasan,
    aa."NAMA" AS a_nama_atasan,
    aa."PNS_ID" AS a_pns_id_atasan,
    ag."ID" AS golongan_atasan,
    a."TMT_GOLONGAN" AS tmt_golongan_atasan,
    aag."ID" AS golongan_atasan_atasan,
    aa."TMT_GOLONGAN" AS tmt_golongan_atasan_atasan,
    a.status_pegawai AS status_pns_atasan,
    aa.status_pegawai AS status_pns_atasan_atasan,
    p."JENIS_JABATAN_ID",
    au."NAMA_UNOR" AS nama_unor_atasan,
    aau."NAMA_UNOR" AS nama_unor_atasan_atasan
   FROM (((((((rwt_prestasi_kerja s
     LEFT JOIN pegawai p ON (((s."PNS_NIP")::text = (p."NIP_BARU")::text)))
     LEFT JOIN pegawai a ON (((s."ATASAN_LANGSUNG_PNS_NIP")::text = (a."NIP_BARU")::text)))
     LEFT JOIN pegawai aa ON (((s."ATASAN_ATASAN_LANGSUNG_PNS_NIP")::text = (aa."NIP_BARU")::text)))
     LEFT JOIN golongan ag ON (((a."GOL_ID")::text = (ag."ID")::text)))
     LEFT JOIN golongan aag ON (((aa."GOL_ID")::text = (aag."ID")::text)))
     LEFT JOIN vw_unit_list au ON (((a."UNOR_ID")::text = (au."ID")::text)))
     LEFT JOIN vw_unit_list aau ON (((aa."UNOR_ID")::text = (aau."ID")::text)));

CREATE VIEW vw_sync_ds AS
 SELECT ld."ID",
    ld."NIK",
    ld."ID_FILE",
    ld."CREATED_DATE",
    ld."STATUS",
    ld."PROSES_CRON",
    btrim((ds.kategori)::text) AS kategori,
    btrim((ds.nip_sk)::text) AS nip_sk,
    ds.nama_pemilik_sk
   FROM (log_ds ld
     LEFT JOIN tbl_file_ds ds ON (((ld."ID_FILE")::text = (ds.id_file)::text)))
  WHERE ((ld."STATUS" = 2) AND (ds.nip_sk IS NOT NULL) AND (ds.telah_kirim = 1) AND (ds.ds_ok = 1) AND (ld."PROSES_CRON" = 0) AND (ds.is_signed = 1));

CREATE VIEW vw_tte_trx_draft_sk_to_sk AS
 SELECT ttdsk.id,
    ttdsk.id_master_proses,
    ttdsk.nip_sk,
    ttdsk.penandatangan_sk,
    ttdsk.tgl_sk,
    ttdsk.nomor_sk,
    ttdsk.file_template,
    ttdsk.base64pdf_hasil,
    ttdsk.created_date,
    ttdsk.created_by,
    ttdsk.updated_date,
    ttdsk.updated_by,
    ttdsk.id_file,
    ttdsk.tmt_sk,
    ttdsk.nama_pemilik_sk,
    ttdsk.halaman_ttd,
    ttdsk.show_qrcode,
    ttdsk.letak_ttd,
    tfd.is_signed
   FROM (tte_trx_draft_sk ttdsk
     JOIN tbl_file_ds tfd ON (((ttdsk.id_file)::text = (tfd.id_file)::text)));

CREATE VIEW vw_unit_list_asli AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."ESELON_1",
    uk."ESELON_2",
    uk."ESELON_3",
    uk."ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    btrim(concat(es1."NAMA_UNOR", '-', es2."NAMA_UNOR", '-', es3."NAMA_UNOR", '-', es4."NAMA_UNOR"), '-'::text) AS "NAMA_UNOR_FULL"
   FROM ((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)));

CREATE MATERIALIZED VIEW vw_unit_list_bak AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(32) AS "ESELON_1",
    (x.eselon[2])::character varying(32) AS "ESELON_2",
    (x.eselon[3])::character varying(32) AS "ESELON_3",
    (x.eselon[4])::character varying(32) AS "ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WHERE (uk."EXPIRED_DATE" IS NULL)
  WITH NO DATA;

CREATE MATERIALIZED VIEW vw_unit_list_bak2 AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(32) AS "ESELON_1",
    (x.eselon[2])::character varying(32) AS "ESELON_2",
    (x.eselon[3])::character varying(32) AS "ESELON_3",
    (x.eselon[4])::character varying(32) AS "ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WHERE (uk."EXPIRED_DATE" IS NULL)
  WITH NO DATA;

CREATE VIEW vw_unit_list_new AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(36) AS "ESELON_1",
    (x.eselon[2])::character varying(36) AS "ESELON_2",
    (x.eselon[3])::character varying(36) AS "ESELON_3",
    (x.eselon[4])::character varying(36) AS "ESELON_4",
    uk."JENIS_SATKER",
    es1."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1",
    es2."NAMA_UNOR" AS "NAMA_UNOR_ESELON_2",
    es3."NAMA_UNOR" AS "NAMA_UNOR_ESELON_3",
    es4."NAMA_UNOR" AS "NAMA_UNOR_ESELON_4",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WHERE (uk."EXPIRED_DATE" IS NULL);

CREATE VIEW vw_unit_list_pejabat AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    uk."PERATURAN",
    (x.eselon[1])::character varying(32) AS "ESELON_1",
    (x.eselon[2])::character varying(32) AS "ESELON_2",
    (x.eselon[3])::character varying(32) AS "ESELON_3",
    (x.eselon[4])::character varying(32) AS "ESELON_4",
    uk."JENIS_SATKER",
    x."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN",
    p."NIP_BARU",
    p."GELAR_DEPAN",
    p."NAMA" AS "PEJABAT_NAMA",
    p."GELAR_BELAKANG"
   FROM ((unitkerja uk
     LEFT JOIN pegawai p ON (((p."PNS_ID")::text = (uk."PEMIMPIN_PNS_ID")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    (((a."NAMA_UNOR")::text || ' - '::text) || r_1."NAMA_UNOR"),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            r."NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)));

CREATE MATERIALIZED VIEW vw_unit_list_penyajian_data AS
 SELECT uk."NO",
    uk."KODE_INTERNAL",
    uk."ID",
    uk."NAMA_UNOR",
    uk."ESELON_ID",
    uk."CEPAT_KODE",
    uk."NAMA_JABATAN",
    uk."NAMA_PEJABAT",
    uk."DIATASAN_ID",
    uk."INSTANSI_ID",
    uk."PEMIMPIN_NON_PNS_ID",
    uk."PEMIMPIN_PNS_ID",
    uk."JENIS_UNOR_ID",
    uk."UNOR_INDUK",
    uk."JUMLAH_IDEAL_STAFF",
    uk."ORDER",
    uk.deleted,
    uk."IS_SATKER",
    uk."EXPIRED_DATE",
    (x.eselon[1])::character varying(500) AS "ESELON_1",
    (x.eselon[2])::character varying(500) AS "ESELON_2",
    (x.eselon[3])::character varying(500) AS "ESELON_3",
    (x.eselon[4])::character varying(500) AS "ESELON_4",
    uk."JENIS_SATKER",
    (x."NAMA_UNOR"[1])::character varying(500) AS "NAMA_UNOR_ESELON_1",
    (x."NAMA_UNOR"[2])::character varying(500) AS "NAMA_UNOR_ESELON_2",
    (x."NAMA_UNOR"[3])::character varying(500) AS "NAMA_UNOR_ESELON_3",
    (x."NAMA_UNOR"[4])::character varying(500) AS "NAMA_UNOR_ESELON_4",
    uk."NAMA_UNOR" AS "NAMA_UNOR_FULL",
    uk."UNOR_INDUK_PENYETARAAN",
    uk."ABBREVIATION"
   FROM (((((unitkerja uk
     LEFT JOIN unitkerja es1 ON (((es1."ID")::text = (uk."ESELON_1")::text)))
     LEFT JOIN unitkerja es2 ON (((es2."ID")::text = (uk."ESELON_2")::text)))
     LEFT JOIN unitkerja es3 ON (((es3."ID")::text = (uk."ESELON_3")::text)))
     LEFT JOIN unitkerja es4 ON (((es4."ID")::text = (uk."ESELON_4")::text)))
     LEFT JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    (unitkerja."NAMA_UNOR")::text AS "NAMA_UNOR",
                    (unitkerja."ID")::text AS arr_id
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a."ID",
                    ((r_1."NAMA_UNOR" || '#'::text) || (a."NAMA_UNOR")::text),
                    ((r_1.arr_id || '#'::text) || (a."ID")::text)
                   FROM (unitkerja a
                     JOIN r r_1 ON (((r_1."ID")::text = (a."DIATASAN_ID")::text)))
                )
         SELECT r."ID",
            string_to_array(r."NAMA_UNOR", '#'::text) AS "NAMA_UNOR",
            string_to_array(r.arr_id, '#'::text) AS eselon
           FROM r) x ON (((uk."ID")::text = (x."ID")::text)))
  WHERE (uk."EXPIRED_DATE" IS NULL)
  WITH NO DATA;

CREATE VIEW vw_unor_satker_copy1 AS
 SELECT a."ID" AS "ID_UNOR",
    b."ID" AS "ID_SATKER",
    a."NAMA_UNOR",
    b."NAMA_UNOR" AS "NAMA_SATKER",
    c."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
   FROM ((unitkerja a
     JOIN unitkerja b ON (((b."ID")::text = (a."UNOR_INDUK")::text)))
     JOIN unitkerja c ON (((a."ESELON_1")::text = (c."ID")::text)))
  WHERE ((a."UNOR_INDUK")::text IN ( SELECT unitkerja."ID"
           FROM unitkerja unitkerja
          WHERE (unitkerja."IS_SATKER" = (1)::smallint)))
UNION ALL
 SELECT a."ID" AS "ID_UNOR",
    a."ID" AS "ID_SATKER",
    a."NAMA_UNOR",
    a."NAMA_UNOR" AS "NAMA_SATKER",
    b."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
   FROM (unitkerja a
     JOIN unitkerja b ON (((b."ID")::text = (a."UNOR_INDUK")::text)))
  WHERE (a."IS_SATKER" = (1)::smallint);

COMMENT ON VIEW vw_unor_satker_copy1 IS 'Untuk Melihat Daftar Unit Kerja Berdasarkan Satkernya';

CREATE VIEW vw_unor_satker_only_satker AS
 SELECT a."ID" AS "ID_UNOR",
    a."UNOR_INDUK" AS "ID_SATKER",
    a."NAMA_UNOR",
    b."NAMA_UNOR" AS "NAMA_SATKER",
    c."NAMA_UNOR_ESELON_1",
    a."EXPIRED_DATE",
    c.id_eselon_1 AS "ID_ESELON_1"
   FROM ((unitkerja a
     JOIN unitkerja b ON (((a."UNOR_INDUK")::text = (b."ID")::text)))
     JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    unitkerja."ID" AS id_eselon_1,
                    unitkerja."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a_1."ID",
                    r_1.id_eselon_1,
                    r_1."NAMA_UNOR_ESELON_1"
                   FROM (unitkerja a_1
                     JOIN r r_1 ON (((a_1."DIATASAN_ID")::text = (r_1."ID")::text)))
                )
         SELECT r."ID",
            r.id_eselon_1,
            r."NAMA_UNOR_ESELON_1"
           FROM r) c ON (((a."ID")::text = (c."ID")::text)))
  WHERE ((a."IS_SATKER" = 1) AND (a."EXPIRED_DATE" IS NULL));

CREATE VIEW vw_unor_satker_satyalencana AS
 SELECT a."ID" AS "ID_UNOR",
    a."UNOR_INDUK" AS "ID_SATKER",
    a."NAMA_UNOR",
    b."NAMA_UNOR" AS "NAMA_SATKER",
    c."NAMA_UNOR_ESELON_1",
    a."EXPIRED_DATE",
    c.id_eselon_1 AS "ID_ESELON_1"
   FROM ((unitkerja a
     JOIN unitkerja b ON (((a."UNOR_INDUK")::text = (b."ID")::text)))
     JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    unitkerja."ID" AS id_eselon_1,
                    unitkerja."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a_1."ID",
                    r_1.id_eselon_1,
                    r_1."NAMA_UNOR_ESELON_1"
                   FROM (unitkerja a_1
                     JOIN r r_1 ON (((a_1."DIATASAN_ID")::text = (r_1."ID")::text)))
                )
         SELECT r."ID",
            r.id_eselon_1,
            r."NAMA_UNOR_ESELON_1"
           FROM r) c ON (((a."ID")::text = (c."ID")::text)));

CREATE VIEW vw_unor_satker_w_eselonid AS
 SELECT a."ID" AS "ID_UNOR",
    a."UNOR_INDUK" AS "ID_SATKER",
    a."NAMA_UNOR",
    b."NAMA_UNOR" AS "NAMA_SATKER",
    c."NAMA_UNOR_ESELON_1",
    a."ESELON_ID",
    btrim((a."NAMA_JABATAN")::text) AS "NAMA_JABATAN",
    a."DIATASAN_ID",
    c."ID" AS "ESELON_1_ID"
   FROM ((unitkerja a
     JOIN unitkerja b ON (((a."UNOR_INDUK")::text = (b."ID")::text)))
     JOIN ( WITH RECURSIVE r AS (
                 SELECT unitkerja."ID",
                    unitkerja."ID" AS id_eselon_1,
                    unitkerja."NAMA_UNOR" AS "NAMA_UNOR_ESELON_1"
                   FROM unitkerja
                  WHERE ((unitkerja."DIATASAN_ID")::text = 'A8ACA7397AEB3912E040640A040269BB'::text)
                UNION ALL
                 SELECT a_1."ID",
                    r_1.id_eselon_1,
                    r_1."NAMA_UNOR_ESELON_1"
                   FROM (unitkerja a_1
                     JOIN r r_1 ON (((a_1."DIATASAN_ID")::text = (r_1."ID")::text)))
                )
         SELECT r."ID",
            r.id_eselon_1,
            r."NAMA_UNOR_ESELON_1"
           FROM r) c ON (((a."ID")::text = (c."ID")::text)));

CREATE VIEW vw_unor_satker_w_id_eselon1 AS
 SELECT "ID" AS "ID_UNOR",
    "NAMA_UNOR",
    "ABBREVIATION",
    "EXPIRED_DATE",
        CASE
            WHEN (btrim(("NAMA_UNOR")::text) = 'KEMENTERIAN PENDIDIKAN dan KEBUDAYAAN'::text) THEN "ID"
            WHEN (("ESELON_2" IS NULL) OR (btrim(("ESELON_2")::text) = ''::text)) THEN "ESELON_1"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'universitas_dikti'::text) THEN "ESELON_2"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'Politeknik Vokasi'::text) THEN "ESELON_2"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Multimedia Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Televisi Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Radio Pendidikan dan Kebudayaan'::text)) THEN "ESELON_3"
            ELSE "ESELON_2"
        END AS "ID_SATKER",
        CASE
            WHEN (btrim(("NAMA_UNOR")::text) = 'KEMENTERIAN PENDIDIKAN dan KEBUDAYAAN'::text) THEN "NAMA_UNOR"
            WHEN (("ESELON_2" IS NULL) OR (btrim(("ESELON_2")::text) = ''::text)) THEN "NAMA_UNOR_ESELON_1"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'universitas_dikti'::text) THEN "NAMA_UNOR_ESELON_2"
            WHEN (btrim(("NAMA_UNOR_ESELON_1")::text) = 'Politeknik Vokasi'::text) THEN "NAMA_UNOR_ESELON_2"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Multimedia Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Televisi Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            WHEN ((btrim(("NAMA_UNOR_ESELON_1")::text) = 'Sekretariat Jenderal'::text) AND (btrim(("NAMA_UNOR_ESELON_2")::text) = 'Pusat Data dan Teknologi Informasi'::text) AND (btrim(("NAMA_UNOR_ESELON_3")::text) = 'Balai Pengembangan Media Radio Pendidikan dan Kebudayaan'::text)) THEN "NAMA_UNOR_ESELON_3"
            ELSE "NAMA_UNOR_ESELON_2"
        END AS "NAMA_SATKER",
        CASE
            WHEN (btrim(("NAMA_UNOR")::text) = 'KEMENTERIAN PENDIDIKAN dan KEBUDAYAAN'::text) THEN "ID"
            ELSE "ESELON_1"
        END AS "ID_ESELON_1",
        CASE
            WHEN (btrim(("NAMA_UNOR")::text) = 'KEMENTERIAN PENDIDIKAN dan KEBUDAYAAN'::text) THEN "NAMA_UNOR"
            ELSE "NAMA_UNOR_ESELON_1"
        END AS "NAMA_UNOR_ESELON_1"
   FROM vw_unit_list_penyajian_data t;

CREATE TABLE wage (
    "GOLONGAN" character varying(2) NOT NULL,
    "WORKING_PERIOD" smallint NOT NULL,
    "BASIC" integer,
    "TSP" integer,
    "ETC" integer
);

ALTER TABLE ONLY absen ALTER COLUMN "ID" SET DEFAULT nextval('"absen_ID_seq"'::regclass);

ALTER TABLE ONLY anak ALTER COLUMN "ID" SET DEFAULT nextval('"anak_ID_seq"'::regclass);

ALTER TABLE ONLY arsip ALTER COLUMN "ID" SET DEFAULT nextval('"arsip_ID_seq"'::regclass);

ALTER TABLE ONLY asesmen_hasil_asesmen ALTER COLUMN id SET DEFAULT nextval('asesmen_hasil_asesmen_id_seq'::regclass);

ALTER TABLE ONLY asesmen_pegawai_berpotensi_jpt ALTER COLUMN id SET DEFAULT nextval('asesmen_pegawai_berpotensi_jpt_id_seq'::regclass);

ALTER TABLE ONLY asesmen_riwayat_hukuman_disiplin ALTER COLUMN id SET DEFAULT nextval('asesmen_riwayat_hukuman_disiplin_id_seq'::regclass);

ALTER TABLE ONLY baperjakat ALTER COLUMN "ID" SET DEFAULT nextval('"baperjakat_ID_seq"'::regclass);

ALTER TABLE ONLY daftar_rohaniawan ALTER COLUMN id SET DEFAULT nextval('daftar_rohaniawan_id_seq'::regclass);

ALTER TABLE ONLY hari_libur ALTER COLUMN "ID" SET DEFAULT nextval('"hari_libur_ID_seq"'::regclass);

ALTER TABLE ONLY istri ALTER COLUMN "ID" SET DEFAULT nextval('"istri_ID_seq"'::regclass);

ALTER TABLE ONLY izin ALTER COLUMN "ID" SET DEFAULT nextval('"izin_ID_seq"'::regclass);

ALTER TABLE ONLY izin_alasan ALTER COLUMN "ID" SET DEFAULT nextval('"izin_alasan_ID_seq"'::regclass);

ALTER TABLE ONLY izin_verifikasi ALTER COLUMN "ID" SET DEFAULT nextval('"izin_verifikasi_ID_seq"'::regclass);

ALTER TABLE ONLY jabatan ALTER COLUMN id SET DEFAULT nextval('jabatan_id_seq'::regclass);

ALTER TABLE ONLY jenis_arsip ALTER COLUMN "ID" SET DEFAULT nextval('jenis_arsip_id_seq'::regclass);

ALTER TABLE ONLY jenis_diklat ALTER COLUMN id SET DEFAULT nextval('jenis_diklat_id_seq'::regclass);

ALTER TABLE ONLY jenis_izin ALTER COLUMN "ID" SET DEFAULT nextval('"jenis_izin_ID_seq"'::regclass);

ALTER TABLE ONLY jenis_kursus ALTER COLUMN id SET DEFAULT nextval('jenis_kursus_id_seq'::regclass);

ALTER TABLE ONLY kandidat_baperjakat ALTER COLUMN "ID" SET DEFAULT nextval('"kandidat_baperjakat_ID_seq"'::regclass);

ALTER TABLE ONLY kategori_ds ALTER COLUMN id SET DEFAULT nextval('kategori_ds_id_seq'::regclass);

ALTER TABLE ONLY kategori_jenis_arsip ALTER COLUMN "ID" SET DEFAULT nextval('"kategori_jenis_arsip_ID_seq"'::regclass);

ALTER TABLE ONLY layanan ALTER COLUMN id SET DEFAULT nextval('layanan_id_seq1'::regclass);

ALTER TABLE ONLY layanan_tipe ALTER COLUMN id SET DEFAULT nextval('layanan_id_seq'::regclass);

ALTER TABLE ONLY layanan_usulan ALTER COLUMN id SET DEFAULT nextval('layanan_usulan_id_seq'::regclass);

ALTER TABLE ONLY line_approval_izin ALTER COLUMN "ID" SET DEFAULT nextval('"line_approval_izin_ID_seq"'::regclass);

ALTER TABLE ONLY log_ds ALTER COLUMN "ID" SET DEFAULT nextval('"log_ds_ID_seq"'::regclass);

ALTER TABLE ONLY log_request ALTER COLUMN id SET DEFAULT nextval('log_request_id_seq'::regclass);

ALTER TABLE ONLY log_transaksi ALTER COLUMN "ID" SET DEFAULT nextval('"log_transaksi_ID_seq"'::regclass);

ALTER TABLE ONLY mst_jenis_satker ALTER COLUMN id_jenis SET DEFAULT nextval('jenis_satker_id_jenis_seq'::regclass);

ALTER TABLE ONLY mst_peraturan_otk ALTER COLUMN id_peraturan SET DEFAULT nextval('peraturan_otk_id_peraturan_seq'::regclass);

ALTER TABLE ONLY mst_templates ALTER COLUMN id SET DEFAULT nextval('mst_templates_id_seq'::regclass);

ALTER TABLE ONLY nip_pejabat ALTER COLUMN id SET DEFAULT nextval('nip_pejabat_id_seq'::regclass);

ALTER TABLE ONLY orang_tua ALTER COLUMN "ID" SET DEFAULT nextval('"orang_tua_ID_seq"'::regclass);

ALTER TABLE ONLY pegawai_atasan ALTER COLUMN "ID" SET DEFAULT nextval('"pegawai_atasan_ID_seq"'::regclass);

ALTER TABLE ONLY pengajuan_tubel ALTER COLUMN "ID" SET DEFAULT nextval('"pengajuan_tubel_ID_seq"'::regclass);

ALTER TABLE ONLY perkiraan_kpo ALTER COLUMN id SET DEFAULT nextval('perkiraan_kpo_id_seq'::regclass);

ALTER TABLE ONLY perkiraan_usulan_log ALTER COLUMN id SET DEFAULT nextval('perkiraan_usulan_log_id_seq'::regclass);

ALTER TABLE ONLY peta_jabatan_permen ALTER COLUMN id SET DEFAULT nextval('peta_jabatan_permen_id_seq'::regclass);

ALTER TABLE ONLY pindah_unit ALTER COLUMN "ID" SET DEFAULT nextval('"pindah_unit_ID_seq"'::regclass);

ALTER TABLE ONLY ref_tunjangan_kinerja ALTER COLUMN "ID" SET DEFAULT nextval('"ref_tunjangan_kinerja_ID_seq"'::regclass);

ALTER TABLE ONLY request_formasi ALTER COLUMN id SET DEFAULT nextval('request_formasi_id_seq'::regclass);

ALTER TABLE ONLY role_permissions ALTER COLUMN id SET DEFAULT nextval('role_permissions_id_seq'::regclass);

ALTER TABLE ONLY roles_users ALTER COLUMN role_user_id SET DEFAULT nextval('roles_users_role_user_id_seq'::regclass);

ALTER TABLE ONLY rpt_golongan_bulan ALTER COLUMN "ID" SET DEFAULT nextval('"rpt_golongan_bulan_ID_seq"'::regclass);

ALTER TABLE ONLY rpt_jumlah_asn ALTER COLUMN "ID" SET DEFAULT nextval('"rpt_jumlah_asn_ID_seq"'::regclass);

ALTER TABLE ONLY rpt_pendidikan_bulan ALTER COLUMN "ID" SET DEFAULT nextval('"rpt_pendidikan_bulan_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_assesmen ALTER COLUMN "ID" SET DEFAULT nextval('"rwt_assesmen_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_diklat ALTER COLUMN id SET DEFAULT nextval('rwt_diklat_id_seq'::regclass);

ALTER TABLE ONLY rwt_hukdis ALTER COLUMN "ID" SET DEFAULT nextval('"rwt_hukdis_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_kgb ALTER COLUMN id SET DEFAULT nextval('rwt_kgb_id_seq'::regclass);

ALTER TABLE ONLY rwt_kinerja ALTER COLUMN id SET DEFAULT nextval('rwt_kinerja_id_seq'::regclass);

ALTER TABLE ONLY rwt_kursus ALTER COLUMN "ID" SET DEFAULT nextval('"rwt_kursus_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_nine_box ALTER COLUMN "ID" SET DEFAULT nextval('"NINE_BOX_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_penghargaan_umum ALTER COLUMN id SET DEFAULT nextval('rwt_penghargaan_umum_id_seq'::regclass);

ALTER TABLE ONLY rwt_penugasan ALTER COLUMN id SET DEFAULT nextval('rwt_penugasan_id_seq'::regclass);

ALTER TABLE ONLY rwt_pns_cpns ALTER COLUMN "ID" SET DEFAULT nextval('"rwt_pns_cpns_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_tugas_belajar ALTER COLUMN "ID" SET DEFAULT nextval('"rwt_tugas_belajar_ID_seq"'::regclass);

ALTER TABLE ONLY rwt_ujikom ALTER COLUMN id SET DEFAULT nextval('rwt_ujikom_id_seq'::regclass);

ALTER TABLE ONLY settings ALTER COLUMN id SET DEFAULT nextval('settings_id_seq'::regclass);

ALTER TABLE ONLY sisa_cuti ALTER COLUMN "ID" SET DEFAULT nextval('"sisa_cuti_ID_seq"'::regclass);

ALTER TABLE ONLY synch_jumlah_pegawai ALTER COLUMN id SET DEFAULT nextval('synch_jumlah_pegawai_id_seq'::regclass);

ALTER TABLE ONLY tbl_file_ds ALTER COLUMN id SET DEFAULT nextval('tbl_file_ds_id_seq'::regclass);

ALTER TABLE ONLY tbl_file_ds_corrector ALTER COLUMN id SET DEFAULT nextval('tbl_file_ds_corrector_id_seq'::regclass);

ALTER TABLE ONLY tbl_file_ds_riwayat ALTER COLUMN id_riwayat SET DEFAULT nextval('tbl_file_ds_riwayat_id_riwayat_seq'::regclass);

ALTER TABLE ONLY tte_master_korektor ALTER COLUMN id SET DEFAULT nextval('tte_master_korektor_id_seq'::regclass);

ALTER TABLE ONLY tte_master_proses ALTER COLUMN id SET DEFAULT nextval('tte_master_proses_id_seq'::regclass);

ALTER TABLE ONLY tte_master_proses_variable ALTER COLUMN id SET DEFAULT nextval('tte_master_proses_variable_id_seq'::regclass);

ALTER TABLE ONLY tte_master_variable ALTER COLUMN id SET DEFAULT nextval('"tte_ master_variable_id_seq"'::regclass);

ALTER TABLE ONLY tte_trx_draft_sk ALTER COLUMN id SET DEFAULT nextval('tte_trx_draft_sk_id_seq'::regclass);

ALTER TABLE ONLY tte_trx_draft_sk_detil ALTER COLUMN id SET DEFAULT nextval('tte_trx_draft_sk_detil_id_seq'::regclass);

ALTER TABLE ONLY tte_trx_korektor_draft ALTER COLUMN id SET DEFAULT nextval('tte_trx_korektor_draft_id_seq'::regclass);

ALTER TABLE ONLY update_mandiri ALTER COLUMN "ID" SET DEFAULT nextval('"update_mandiri_ID_seq"'::regclass);

ALTER TABLE ONLY usulan_dokumen ALTER COLUMN id SET DEFAULT nextval('usulan_documents_id_seq'::regclass);

ALTER TABLE ONLY rwt_nine_box
    ADD CONSTRAINT "NINE_BOX_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY absen
    ADD CONSTRAINT absen_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (activity_id);

ALTER TABLE ONLY agama
    ADD CONSTRAINT agama_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY anak
    ADD CONSTRAINT anak_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY arsip
    ADD CONSTRAINT arsip_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY asesmen_hasil_asesmen
    ADD CONSTRAINT asesmen_hasil_asesmen_pkey PRIMARY KEY (id);

ALTER TABLE ONLY asesmen_pegawai_berpotensi_jpt
    ADD CONSTRAINT asesmen_pegawai_berpotensi_jpt_pkey PRIMARY KEY (id);

ALTER TABLE ONLY asesmen_riwayat_hukuman_disiplin
    ADD CONSTRAINT asesmen_riwayat_hukuman_disiplin_pkey PRIMARY KEY (id);

ALTER TABLE ONLY ref_tunjangan_jabatan
    ADD CONSTRAINT data_jabatan_tunjab_pkey PRIMARY KEY ("ID_TUNJAB");

ALTER TABLE ONLY golongan
    ADD CONSTRAINT golongan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY hari_libur
    ADD CONSTRAINT hari_libur_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY instansi
    ADD CONSTRAINT instansi_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY istri
    ADD CONSTRAINT istri_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY izin_alasan
    ADD CONSTRAINT izin_alasan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY izin_verifikasi
    ADD CONSTRAINT izin_verifikasi_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jabatan
    ADD CONSTRAINT jabatan_pkey PRIMARY KEY ("KODE_JABATAN");

ALTER TABLE ONLY jenis_arsip
    ADD CONSTRAINT jenis_arsip_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_diklat_fungsional
    ADD CONSTRAINT jenis_diklat_fungsional_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_diklat
    ADD CONSTRAINT jenis_diklat_pkey PRIMARY KEY (id);

ALTER TABLE ONLY jenis_diklat_siasn
    ADD CONSTRAINT jenis_diklat_siasn_pkey PRIMARY KEY (id);

ALTER TABLE ONLY jenis_diklat_struktural
    ADD CONSTRAINT jenis_diklat_struktural_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_jabatan
    ADD CONSTRAINT jenis_jabatan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_kawin
    ADD CONSTRAINT jenis_kawin_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_kp
    ADD CONSTRAINT jenis_kp_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_kursus
    ADD CONSTRAINT jenis_kursus_pkey PRIMARY KEY (id);

ALTER TABLE ONLY jenis_pegawai
    ADD CONSTRAINT jenis_pegawai_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_penghargaan
    ADD CONSTRAINT "jenis_penghargaan_ID" PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_penghargaan
    ADD CONSTRAINT "jenis_penghargaan_NAMA" UNIQUE ("NAMA");

ALTER TABLE ONLY jenis_rumpun_diklat_siasn
    ADD CONSTRAINT jenis_rumpun_diklat_siasn_pkey PRIMARY KEY (id);

ALTER TABLE ONLY mst_jenis_satker
    ADD CONSTRAINT jenis_satker_pkey PRIMARY KEY (id_jenis);

ALTER TABLE ONLY kandidat_baperjakat
    ADD CONSTRAINT kandidat_baperjakat_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY kategori_ds
    ADD CONSTRAINT kategori_ds_pkey PRIMARY KEY (id);

ALTER TABLE ONLY kategori_jenis_arsip
    ADD CONSTRAINT kategori_jenis_arsip_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY kedudukan_hukum
    ADD CONSTRAINT kedudukan_hukum_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY kpkn
    ADD CONSTRAINT kpkn_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY layanan_tipe
    ADD CONSTRAINT layanan_pkey PRIMARY KEY (id);

ALTER TABLE ONLY layanan
    ADD CONSTRAINT layanan_pkey1 PRIMARY KEY (id);

ALTER TABLE ONLY layanan_usulan
    ADD CONSTRAINT layanan_usulan_pkey PRIMARY KEY (id);

ALTER TABLE ONLY line_approval_izin
    ADD CONSTRAINT line_approval_izin_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY log_ds
    ADD CONSTRAINT log_ds_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY log_transaksi
    ADD CONSTRAINT log_transaksi_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY login_attempts
    ADD CONSTRAINT login_attempts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY lokasi
    ADD CONSTRAINT lokasi_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY mst_templates
    ADD CONSTRAINT mst_templates_pkey PRIMARY KEY (id);

ALTER TABLE ONLY orang_tua
    ADD CONSTRAINT orang_tua_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY pegawai_atasan
    ADD CONSTRAINT pegawai_atasan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY pegawai_bkn
    ADD CONSTRAINT pegawai_bkn_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY pegawai
    ADD CONSTRAINT pegawai_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY pendidikan
    ADD CONSTRAINT pendidikan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY mst_peraturan_otk
    ADD CONSTRAINT peraturan_otk_pkey PRIMARY KEY (id_peraturan);

ALTER TABLE ONLY perkiraan_ppo
    ADD CONSTRAINT perkiraan_kpo_copy1_pkey1 PRIMARY KEY (id);

ALTER TABLE ONLY usulan_dokumen
    ADD CONSTRAINT perkiraan_kpo_documents_pkey PRIMARY KEY (id);

ALTER TABLE ONLY perkiraan_kpo
    ADD CONSTRAINT perkiraan_kpo_pkey PRIMARY KEY (id);

ALTER TABLE ONLY permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (permission_id);

ALTER TABLE ONLY peta_jabatan_permen
    ADD CONSTRAINT peta_jabatan_permen_pkey PRIMARY KEY (id);

ALTER TABLE ONLY baperjakat
    ADD CONSTRAINT pk_baperjakat PRIMARY KEY ("ID");

ALTER TABLE ONLY daftar_rohaniawan
    ADD CONSTRAINT pk_daftar_rohaniawan PRIMARY KEY (id);

ALTER TABLE ONLY izin
    ADD CONSTRAINT pk_izin PRIMARY KEY ("ID");

ALTER TABLE ONLY jenis_izin
    ADD CONSTRAINT pk_jenis_izin PRIMARY KEY ("ID");

ALTER TABLE ONLY pengajuan_tubel
    ADD CONSTRAINT pk_pengajuan_tubel PRIMARY KEY ("ID");

ALTER TABLE ONLY pindah_unit
    ADD CONSTRAINT pk_pindah_unit PRIMARY KEY ("ID");

ALTER TABLE ONLY sisa_cuti
    ADD CONSTRAINT pk_sisa_cuti PRIMARY KEY ("ID");

ALTER TABLE ONLY tte_master_proses
    ADD CONSTRAINT pk_tte_master_proses PRIMARY KEY (id);

ALTER TABLE ONLY ref_jabatan
    ADD CONSTRAINT ref_jabatan_pkey PRIMARY KEY ("ID_JABATAN");

ALTER TABLE ONLY request_formasi
    ADD CONSTRAINT request_formasi_pkey PRIMARY KEY (id);

ALTER TABLE ONLY role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (role_id);

ALTER TABLE ONLY rpt_golongan_bulan
    ADD CONSTRAINT rpt_golongan_bulan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rpt_jumlah_asn
    ADD CONSTRAINT rpt_jumlah_asn_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rpt_pendidikan_bulan
    ADD CONSTRAINT rpt_pendidikan_bulan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_assesmen
    ADD CONSTRAINT rwt_assesmen_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_diklat_fungsional
    ADD CONSTRAINT rwt_diklat_fungsional_pkey PRIMARY KEY ("DIKLAT_FUNGSIONAL_ID");

ALTER TABLE ONLY rwt_diklat
    ADD CONSTRAINT rwt_diklat_pkey PRIMARY KEY (id);

ALTER TABLE ONLY rwt_diklat_struktural
    ADD CONSTRAINT rwt_diklat_struktural_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_golongan
    ADD CONSTRAINT rwt_golongan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_jabatan_empty
    ADD CONSTRAINT rwt_jabatan_copy1_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_jabatan
    ADD CONSTRAINT rwt_jabatan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_pekerjaan
    ADD CONSTRAINT rwt_pekerjaan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_pendidikan
    ADD CONSTRAINT rwt_pendidikan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_penghargaan
    ADD CONSTRAINT "rwt_penghargaan_ID" PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_penghargaan_umum
    ADD CONSTRAINT rwt_penghargaan_umum_pkey PRIMARY KEY (id);

ALTER TABLE ONLY rwt_penugasan
    ADD CONSTRAINT rwt_penugasan_pkey PRIMARY KEY (id);

ALTER TABLE ONLY rwt_pindah_unit_kerja
    ADD CONSTRAINT rwt_pindah_unit_kerja_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_pns_cpns
    ADD CONSTRAINT rwt_pns_cpns_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_prestasi_kerja
    ADD CONSTRAINT rwt_prestasi_kerja_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_tugas_belajar
    ADD CONSTRAINT rwt_tugas_belajar_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY rwt_ujikom
    ADD CONSTRAINT rwt_ujikom_pkey PRIMARY KEY (id);

ALTER TABLE ONLY settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);

ALTER TABLE ONLY synch_jumlah_pegawai
    ADD CONSTRAINT synch_jumlah_pegawai_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tb_nomor_surat
    ADD CONSTRAINT tb_nomor_surat_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tbl_file_ds_corrector
    ADD CONSTRAINT tbl_file_ds_corrector_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tbl_file_ds_khusus_login
    ADD CONSTRAINT tbl_file_ds_khusus_login_pkey PRIMARY KEY ("ID_FILE");

ALTER TABLE ONLY tbl_file_ds
    ADD CONSTRAINT tbl_file_ds_pkey PRIMARY KEY (id_file);

ALTER TABLE ONLY tbl_file_ttd
    ADD CONSTRAINT tbl_file_ttd_pkey PRIMARY KEY (id_pns_bkn);

ALTER TABLE ONLY tbl_kategori_dokumen_penandatangan
    ADD CONSTRAINT tbl_kategori_dokumen_penandatangan_pkey PRIMARY KEY ("ID_URUT");

ALTER TABLE ONLY tbl_kategori_dokumen
    ADD CONSTRAINT tbl_kategori_dokumen_pkey PRIMARY KEY (id_kategori);

ALTER TABLE ONLY tbl_pengantar_dokumen
    ADD CONSTRAINT tbl_pengantar_dokumen_pkey PRIMARY KEY (id_pengantar);

ALTER TABLE ONLY tkpendidikan
    ADD CONSTRAINT tkpendidikan_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY tte_master_variable
    ADD CONSTRAINT "tte_ master_variable_pkey" PRIMARY KEY (id);

ALTER TABLE ONLY tte_master_korektor
    ADD CONSTRAINT tte_master_korektor_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tte_master_proses_variable
    ADD CONSTRAINT tte_master_proses_variable_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tte_trx_draft_sk_detil
    ADD CONSTRAINT tte_trx_draft_sk_detil_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tte_trx_draft_sk
    ADD CONSTRAINT tte_trx_draft_sk_pkey PRIMARY KEY (id);

ALTER TABLE ONLY tte_trx_korektor_draft
    ADD CONSTRAINT tte_trx_korektor_draft_pkey PRIMARY KEY (id);

ALTER TABLE ONLY unitkerja
    ADD CONSTRAINT unitkerja_pkey1 PRIMARY KEY ("ID");

ALTER TABLE ONLY update_mandiri
    ADD CONSTRAINT update_mandiri_pkey PRIMARY KEY ("ID");

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY wage
    ADD CONSTRAINT wage_pkey PRIMARY KEY ("GOLONGAN", "WORKING_PERIOD");

CREATE INDEX "pegawai_GOL_ID" ON pegawai USING btree ("GOL_ID");

CREATE UNIQUE INDEX "pegawai_NIP_BARU" ON pegawai USING btree ("NIP_BARU");

CREATE UNIQUE INDEX "pegawai_PNS_ID_idx" ON pegawai USING btree ("PNS_ID");

CREATE INDEX pegawai_unor_id ON pegawai USING btree ("UNOR_ID");

CREATE INDEX rwt_assesmen_nip ON rwt_assesmen USING btree ("PNS_NIP");

CREATE INDEX rwt_assesmen_pns_id ON rwt_assesmen USING btree ("PNS_ID");

CREATE INDEX "rwt_diklat_fungsional_NIP" ON rwt_diklat_fungsional USING btree ("NIP_BARU");

CREATE INDEX "rwt_diklat_struktural_NIP" ON rwt_diklat_struktural USING btree ("PNS_NIP");

CREATE INDEX rwt_diklat_struktural_pns_id ON rwt_diklat_struktural USING btree ("PNS_ID");

CREATE INDEX "rwt_penghargaan_ID_GOLONGAN" ON rwt_penghargaan USING btree ("ID_GOLONGAN");

CREATE INDEX "rwt_penghargaan_ID_PENGHARGAAN" ON rwt_penghargaan USING btree ("ID_JENIS_PENGHARGAAN");

CREATE INDEX "rwt_penghargaan_PNS_ID" ON rwt_penghargaan USING btree ("PNS_ID");

CREATE INDEX "rwt_penghargaan_PNS_NIP" ON rwt_penghargaan USING btree ("PNS_NIP");

CREATE UNIQUE INDEX settings_name_idx ON settings USING btree (name);

CREATE INDEX username_index ON users USING btree (username);
