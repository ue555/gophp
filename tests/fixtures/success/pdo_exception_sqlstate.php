<?php
$pdo = new PDO('sqlite::memory:');
try {
    $pdo->query('SELECT * FROM absent');
} catch (PDOException $e) {
    echo $e->getCode(), '|', $e->errorInfo[0], '|', $e->errorInfo[1], '|', $e->errorInfo[2], '|';
}

$pdo->exec('CREATE TABLE entries (id INTEGER PRIMARY KEY, name TEXT UNIQUE)');
$pdo->exec("INSERT INTO entries (name) VALUES ('same')");
try {
    $pdo->exec("INSERT INTO entries (name) VALUES ('same')");
} catch (PDOException $e) {
    echo $e->getCode(), '|', $e->errorInfo[0], '|', $e->errorInfo[1], '|', $e->errorInfo[2];
}
