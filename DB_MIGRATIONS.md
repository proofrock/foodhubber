# Migrating databases

## `v0.6.3`->`v0.7.0`

```sql
DELETE FROM rules WHERE profile = 'A3';
UPDATE rules SET profile = 'A' WHERE profile = 'A1';
UPDATE rules SET quantity_w1 = quantity_w2, quantity_w2 = quantity_w4, quantity_w4 = 0 WHERE profile = 'B';
UPDATE beneficiaries SET profile = 'A' WHERE profile IN ('A1', 'A3');
DROP VIEW vu_enabled_weeks;
ALTER TABLE rules RENAME COLUMN quantity_w1 TO quantity_o1;
ALTER TABLE rules RENAME COLUMN quantity_w2 TO quantity_o2;
ALTER TABLE rules RENAME COLUMN quantity_w3 TO quantity_o3;
ALTER TABLE rules RENAME COLUMN quantity_w4 TO quantity_o4;
CREATE VIEW vu_monthly_orders_by_profile AS
  WITH SUMS AS (
  SELECT profile, SUM(quantity_o2) as o2,
         SUM(quantity_o3) as o3, SUM(quantity_o4) as o4
    FROM rules 
   GROUP BY profile
  )
  SELECT profile,
         CASE WHEN o2 = 0 THEN 1
              WHEN o3 = 0 THEN 2
              WHEN o4 = 0 THEN 3
              ELSE             4 END AS num
    FROM SUMS;
```
