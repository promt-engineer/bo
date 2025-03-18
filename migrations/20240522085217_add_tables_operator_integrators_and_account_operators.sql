-- +goose Up
-- +goose StatementBegin
-- ----------------------------
-- Table structure for account_operators
-- ----------------------------
DROP TABLE IF EXISTS "public"."account_operators";
CREATE TABLE "public"."account_operators" (
                                              "created_at" timestamptz(6) DEFAULT now(),
                                              "account_id" uuid NOT NULL,
                                              "operator_id" uuid NOT NULL
);

-- ----------------------------
-- Table structure for operator_integrators
-- ----------------------------
DROP TABLE IF EXISTS "public"."operator_integrators";
CREATE TABLE "public"."operator_integrators" (
                                                 "id" uuid NOT NULL,
                                                 "integrator_id" uuid,
                                                 "operator_id" uuid
);

-- ----------------------------
-- Primary Key structure for table account_operators
-- ----------------------------
ALTER TABLE "public"."account_operators" ADD CONSTRAINT "account_operators_pkey" PRIMARY KEY ("account_id", "operator_id");

-- ----------------------------
-- Primary Key structure for table operator_integrators
-- ----------------------------
ALTER TABLE "public"."operator_integrators" ADD CONSTRAINT "operator_integrators_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table account_operators
-- ----------------------------
ALTER TABLE "public"."account_operators" ADD CONSTRAINT "account_operators_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."account_operators" ADD CONSTRAINT "account_operators_operator_id_fkey" FOREIGN KEY ("operator_id") REFERENCES "public"."organizations" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table operator_integrators
-- ----------------------------
ALTER TABLE "public"."operator_integrators" ADD CONSTRAINT "operator_integrators_operator_id_fkey" FOREIGN KEY ("operator_id") REFERENCES "public"."organizations" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."operator_integrators" ADD CONSTRAINT "operator_integrators_integrator_id_fkey" FOREIGN KEY ("integrator_id") REFERENCES "public"."organizations" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."account_operators";
DROP TABLE IF EXISTS "public"."operator_integrators";
-- +goose StatementEnd