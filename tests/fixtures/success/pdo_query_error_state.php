<?php
$pdo = new PDO('sqlite::memory:');
$pdo->exec('CREATE TABLE entries (name TEXT UNIQUE)');
$pdo->query("INSERT INTO entries (name) VALUES ('same')");
try {
    $pdo->query("INSERT INTO entries (name) VALUES ('same')");
} catch (PDOException $e) {
    echo 'query|';
}
echo $pdo->errorCode(), '|', $pdo->errorInfo()[1], '|';
$pdo->exec("INSERT INTO entries (name) VALUES ('fresh')");
echo $pdo->errorCode(), '|', $pdo->errorInfo()[1] === null ? 'clear' : 'stale';
