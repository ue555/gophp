<?php
$pdo = new PDO('sqlite::memory:', null, null, [PDO::ATTR_ERRMODE => PDO::ERRMODE_SILENT]);
$pdo->exec('CREATE TABLE entries (name TEXT UNIQUE)');
$stmt = $pdo->prepare('INSERT INTO entries (name) VALUES (?)');
$stmt->execute(['same']);
echo $stmt->execute(['same']) === false ? 'statement-false|' : 'wrong|';
echo $stmt->errorCode(), '|', $pdo->errorCode(), '|';
echo $pdo->query("INSERT INTO entries (name) VALUES ('same')") === false ? 'query-false|' : 'wrong|';
echo $pdo->errorCode(), '|';
$pdo->exec("INSERT INTO entries (name) VALUES ('fresh')");
echo $pdo->errorCode();
