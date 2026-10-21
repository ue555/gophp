<?php
require 'research/laravel-api/vendor/autoload.php';

$capsule = new Illuminate\Database\Capsule\Manager;
$capsule->addConnection(['driver' => 'sqlite', 'database' => ':memory:', 'prefix' => '']);
$db = $capsule->getConnection();
$direct = $db->cursor('SELECT id,title FROM items WHERE id > ? ORDER BY id', [1]);
$wrapped = $db->table('items')->orderBy('id')->cursor();
echo 'created|';

$db->statement('CREATE TABLE items (id INTEGER PRIMARY KEY, title TEXT)');
$db->insert('INSERT INTO items (id,title) VALUES (?, ?)', [1, 'alpha']);
$db->insert('INSERT INTO items (id,title) VALUES (?, ?)', [2, 'beta']);

foreach ($direct as $row) {
    echo $row->id, ':', $row->title, '|';
}

echo 'wrapped|';
foreach ($wrapped as $row) {
    echo $row->id, ':', $row->title, '|';
}
