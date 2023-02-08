-- шляпа
SELECT dil.id, dil.name, dil.phone_number, json_agg(it) as item
FROM (SELECT i.id, i.name, i.is_legal, di.id
      FROM items i
               JOIN dialer_item di ON i.id = di.item_id) it
         JOIN (SELECT d.id, d.name, d.phone_number, di2.id
               FROM dialers d
                        JOIN dialer_item di2 ON d.id = di2.dialer_id) dil
              ON dil.id = it.id
WHERE it.id = 1
GROUP BY dil.id, dil.name, dil.phone_number;


-- тож шляпа но близко (обосрался в агригации нужно чекнуть сигнатуру)
SELECT dil.dialer_id, dil.dialer_name, dil.phone_number, json_agg(dil) AS items
FROM (SELECT (d.id)   as dialer_id,
             (d.name) as dialer_name,
             d.phone_number,
             (i.id)   as item_id,
             (i.name) as item_name,
             i.is_legal
      FROM dialers d
               JOIN dialer_item di ON d.id = di.dialer_id
               JOIN items i ON i.id = di.item_id
      GROUP BY d.id, d.name, d.phone_number, i.id) dil
WHERE dil.item_id = 2;


-- уже близко
SELECT dialer.id, dialer.name, dialer.phone_number, json_agg(item) AS items
FROM (SELECT d.id, d.name, d.phone_number, (di.id) AS dialer_id
      FROM dialers d
               JOIN dialer_item di ON d.id = di.dialer_id) dialer
         JOIN (SELECT i.id, i.name, i.is_legal, (di2.id) AS item_id
               FROM items i
                        JOIN dialer_item di2 ON i.id = di2.item_id) item
              ON dialer.dialer_id = item.item_id
GROUP BY dialer.id, dialer.name, dialer.phone_number;


SELECT d.id, d.name, d.phone_number, json_agg(i)
FROM dialers d JOIN dialer_item di ON d.id = di.dialer_id
               JOIN items i ON i.id = di.item_id
GROUP BY d.id, d.name, d.phone_number;

-- идеально
select d.id, d.name, d.phone_number, json_agg(i) as items
from dialers d
join dialer_item di on d.id = di.dialer_id
join items i on di.item_id = i.id
where d.id IN
(select d.id from dialers d join dialer_item di on d.id = di.dialer_id where di.item_id = 1)
group by d.id;