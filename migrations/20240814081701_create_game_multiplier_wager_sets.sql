-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."integrator_game_wager_sets";
CREATE TABLE "public"."integrator_game_wager_sets" (
                                                       "organization_id" uuid NOT NULL,
                                                       "game_id" uuid NOT NULL,
                                                       "currency" VARCHAR(255) NOT NULL,
                                                       "wager_set_id" uuid NOT NULL
)
;

-- ----------------------------
-- Primary Key structure for table integrator_games
-- ----------------------------
ALTER TABLE "public"."integrator_game_wager_sets" ADD CONSTRAINT "integrator_game_wager_sets_pkey" PRIMARY KEY ("organization_id", "game_id", "currency", "wager_set_id");

-- ----------------------------
-- Foreign Keys structure for table integrator_games
-- ----------------------------
ALTER TABLE "public"."integrator_game_wager_sets"
    ADD CONSTRAINT "integrator_game_wager_sets_game_id_fkey" FOREIGN KEY ("game_id") REFERENCES "public"."games" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
    ADD CONSTRAINT "integrator_game_wager_sets_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON DELETE CASCADE ON UPDATE NO ACTION,
    ADD CONSTRAINT "integrator_game_wager_sets_wager_set_id_fkey" FOREIGN KEY ("wager_set_id") REFERENCES "public"."wager_sets" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."integrator_game_wager_sets";
-- +goose StatementEnd
