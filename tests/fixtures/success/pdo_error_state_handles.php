<?php
$pdo = new PDO('sqlite::memory:');
echo $pdo->errorCode() === null ? 'fresh|' : 'stale|';
echo $pdo->errorInfo()[0] === '' ? 'empty|' : 'nonempty|';

try {
    $pdo->exec('INSERT INTO absent VALUES (1)');
} catch (PDOException $e) {
    echo 'direct|';
}
echo $pdo->errorCode(), '|', $pdo->errorInfo()[1], '|';

$pdo->exec('CREATE TABLE entries (name TEXT UNIQUE)');
echo $pdo->errorCode(), '|', $pdo->errorInfo()[1] === null ? 'clear|' : 'stale|';

$stmt = $pdo->prepare('INSERT INTO entries (name) VALUES (?)');
echo $stmt->errorCode() === null && $stmt->errorInfo()[0] === '' ? 'unrun|' : 'run|';
$stmt->execute(['same']);
echo $stmt->errorCode(), '|';
try {
    $stmt->execute(['same']);
} catch (PDOException $e) {
    echo 'statement|';
}
echo $stmt->errorCode(), '|', $stmt->errorInfo()[1], '|', $pdo->errorCode(), '|';
$stmt->execute(['fresh']);
echo $stmt->errorCode(), '|', $stmt->errorInfo()[1] === null ? 'clear' : 'stale';
