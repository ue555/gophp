<?php
require 'research/laravel-api/vendor/autoload.php';

$capsule = new Illuminate\Database\Capsule\Manager;
$capsule->addConnection(['driver' => 'sqlite', 'database' => ':memory:', 'prefix' => '']);
$db = $capsule->getConnection();
$db->statement('CREATE TABLE entries (id INTEGER PRIMARY KEY, name TEXT UNIQUE)');
$db->insert('INSERT INTO entries (name) VALUES (?)', ['same']);
try {
    $db->insert('INSERT INTO entries (name) VALUES (?)', ['same']);
} catch (Illuminate\Database\UniqueConstraintViolationException $e) {
    echo $e->getCode(), '|', get_class($e->getPrevious()), '|', $e->errorInfo[0];
}
