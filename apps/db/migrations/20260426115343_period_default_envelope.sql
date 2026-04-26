-- migrate:up
ALTER TABLE financial_periods
  ADD COLUMN default_envelope_id UUID REFERENCES envelopes(id) ON DELETE SET NULL;

-- migrate:down
ALTER TABLE financial_periods DROP COLUMN default_envelope_id;
