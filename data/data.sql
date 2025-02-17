INSERT INTO configs (key, value) VALUES
('order_list_page_size', '25'),
('polling_cycle', '1000'),
('yellow_limit', '150'),
('red_limit', '25');

INSERT INTO items
(id, pos, color, item, subitem)
VALUES
(010, 010, '#FFADAD', 'Olio', 'Oliva'),
(020, 020, '#FFADAD', 'Olio', 'Semi'),
(030, 030, '#FFADAD', 'Riso', null),
(040, 040, '#FFADAD', 'Sale', null),
(050, 050, '#FFADAD', 'Zucchero', null),
(060, 060, '#FFADAD', 'Farina', null),
(070, 070, '#FFADAD', 'Extra Salato', null),
(080, 080, '#FFADAD', 'Extra Dolce', null),
(090, 090, '#FFADAD', 'Crackers', null),
--
(100, 100, '#C6DEF1', 'Latte', null),
(110, 110, '#C6DEF1', 'Caffè', null),
(120, 120, '#C6DEF1', 'Colazione', 'Fette Biscottate'),
(130, 130, '#C6DEF1', 'Colazione', 'Frollini'),
(140, 140, '#C6DEF1', 'Colazione', 'Crostatine'),
(150, 150, '#C6DEF1', 'Spalmabili', 'Confettura'),
(160, 160, '#C6DEF1', 'Spalmabili', 'Crema Cacao'),
--
(170, 170, '#FDFFB6', 'Legumi', 'Ceci'),
(180, 180, '#FDFFB6', 'Legumi', 'Fagioli'),
(190, 190, '#FDFFB6', 'Legumi', 'Lenticchie'),
(200, 200, '#FDFFB6', 'Pasta', 'Spaghetti'),
(210, 210, '#FDFFB6', 'Pasta', 'Stelline'),
(220, 220, '#FDFFB6', 'Pasta', 'Sedani Rigati'),
(230, 230, '#FDFFB6', 'Pasta', 'Penne Rigate'),
(240, 240, '#FDFFB6', 'Pomodoro', 'Polpa'),
(250, 250, '#FDFFB6', 'Pomodoro', 'Passata'),
(260, 260, '#FDFFB6', 'Secondo', 'Tonno'),
(270, 270, '#FDFFB6', 'Secondo', 'Carne in Scatola'),
--
(280, 280, '#D9EDF8', 'Omogeneizzati', null),
(290, 290, '#D9EDF8', 'Pannolini', null),
(300, 300, '#D9EDF8', 'Assorbenti', null),
(310, 310, '#D9EDF8', 'Igiene Persona', 'Shampoo'),
(320, 320, '#D9EDF8', 'Igiene Persona', 'Bagnoschiuma'),
(330, 330, '#D9EDF8', 'Igiene Persona', 'Sapone Mani'),
(340, 340, '#D9EDF8', 'Igiene Casa', 'Piatti'),
(350, 350, '#D9EDF8', 'Igiene Casa', 'Pavimenti'),
(360, 360, '#D9EDF8', 'Igiene Casa', 'Lavatrice'),
--
(370, 370, '#D1F7B5', 'Fresco', 'Frutta e Verdura'),
--
(380, 380, '#FFFFFF', 'Freschissimo', null);

INSERT INTO checkouts
(id, pos, can_access_order_list_page, can_delete_orders, can_access_stats_page,
 can_access_stock_page, can_change_stock, can_access_console_page, password_hash)
VALUES
('CASSA 1', 10, true, false, false, true, false, false, null),
('CASSA 2', 20, true, false, false, true, false, false, null),
('MANAGER', 30, true, true,  true,  true, true,  true, 'e52f60ef556767578566ec462fa1b19b78bcd8049ff9629b5f89174d3a497960');

INSERT INTO rules
(profile, item, quantity_o1, quantity_o2, quantity_o3, quantity_o4)
VALUES
('A', 'Olio', 1, 0, 0, 0),
('A', 'Riso', 1, 0, 0, 0),
('A', 'Sale', 1, 0, 0, 0),
('A', 'Zucchero', 1, 0, 0, 0),
('A', 'Farina', 1, 0, 0, 0),
('A', 'Extra salato', 1, 0, 0, 0),
('A', 'Extra dolce', 1, 0, 0, 0),
('A', 'Crackers', 1, 0, 0, 0),
('A', 'Latte', 2, 0, 0, 0),
('A', 'Caffè', 2, 0, 0, 0),
('A', 'Colazione', 2, 0, 0, 0),
('A', 'Spalmabili', 2, 0, 0, 0),
('A', 'Legumi', 3, 0, 0, 0),
('A', 'Pasta', 3, 0, 0, 0),
('A', 'Pomodoro', 3, 0, 0, 0),
('A', 'Secondo', 4, 0, 0, 0),
('A', 'Omogeneizzati', 6, 0, 0, 0),
('A', 'Pannolini', 1, 0, 0, 0),
('A', 'Assorbenti', 1, 0, 0, 0),
('A', 'Igiene Persona', 1, 0, 0, 0),
('A', 'Igiene Casa', 1, 0, 0, 0),
('A', 'Fresco', 3, 0, 0, 0),
('A', 'Freschissimo', 5, 0, 0, 0),
--
('B', 'Olio', 0, 1, 0, 0),
('B', 'Riso', 0, 1, 0, 0),
('B', 'Sale', 0, 1, 0, 0),
('B', 'Zucchero', 0, 1, 0, 0),
('B', 'Farina', 0, 1, 0, 0),
('B', 'Extra salato', 0, 1, 0, 0),
('B', 'Extra dolce', 0, 1, 0, 0),
('B', 'Crackers', 0, 1, 0, 0),
('B', 'Latte', 0, 2, 0, 2),
('B', 'Caffè', 0, 2, 0, 2),
('B', 'Colazione', 0, 2, 0, 2),
('B', 'Spalmabili', 0, 2, 0, 2),
('B', 'Legumi', 0, 3, 0, 3),
('B', 'Pasta', 0, 3, 0, 3),
('B', 'Pomodoro', 0, 3, 0, 3),
('B', 'Secondo', 0, 4, 0, 4),
('B', 'Omogeneizzati', 0, 6, 0, 0),
('B', 'Pannolini', 0, 1, 0, 0),
('B', 'Assorbenti', 0, 1, 0, 0),
('B', 'Igiene Persona', 0, 1, 0, 0),
('B', 'Igiene Casa', 0, 1, 0, 0),
('B', 'Fresco', 0, 3, 0, 3),
('B', 'Freschissimo', 0, 5, 0, 5),
--
('C', 'Olio', 1, 0, 0, 0),
('C', 'Riso', 1, 0, 0, 0),
('C', 'Sale', 1, 0, 0, 0),
('C', 'Zucchero', 1, 0, 0, 0),
('C', 'Farina', 1, 0, 0, 0),
('C', 'Extra salato', 1, 0, 0, 0),
('C', 'Extra dolce', 1, 0, 0, 0),
('C', 'Crackers', 1, 0, 0, 0),
('C', 'Latte', 2, 2, 2, 2),
('C', 'Caffè', 2, 2, 2, 2),
('C', 'Colazione', 2, 2, 2, 2),
('C', 'Spalmabili', 2, 2, 2, 2),
('C', 'Legumi', 3, 3, 3, 3),
('C', 'Pasta', 3, 3, 3, 3),
('C', 'Pomodoro', 3, 3, 3, 3),
('C', 'Secondo', 4, 4, 4, 4),
('C', 'Omogeneizzati', 6, 0, 0, 0),
('C', 'Pannolini', 1, 0, 0, 0),
('C', 'Assorbenti', 1, 0, 0, 0),
('C', 'Igiene Persona', 1, 0, 0, 0),
('C', 'Igiene Casa', 1, 0, 0, 0),
('C', 'Fresco', 3, 3, 3, 3),
('C', 'Freschissimo', 5, 5, 5, 5);

INSERT INTO beneficiaries
(id, profile)
VALUES
('123', 'A');
