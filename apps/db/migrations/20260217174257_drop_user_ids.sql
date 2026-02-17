-- migrate:up
ALTER TABLE ONLY public.envelopes DROP CONSTRAINT fk_envelopes_user;
DROP INDEX idx_envelopes_user_id;
ALTER TABLE ONLY public.envelopes DROP COLUMN user_id;

ALTER TABLE ONLY public.transactions DROP CONSTRAINT fk_transactions_user;
DROP INDEX idx_transactions_user_id;
ALTER TABLE ONLY public.transactions DROP COLUMN user_id;

-- migrate:down
ALTER TABLE ONLY public.envelopes ADD COLUMN user_id uuid;
CREATE INDEX idx_envelopes_user_id ON public.envelopes USING btree (user_id);
ALTER TABLE ONLY public.envelopes ADD CONSTRAINT fk_envelopes_user FOREIGN KEY (user_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.transactions ADD COLUMN user_id uuid;
CREATE INDEX idx_transactions_user_id ON public.transactions USING btree (user_id);
ALTER TABLE ONLY public.transactions ADD CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES public.users(id);
