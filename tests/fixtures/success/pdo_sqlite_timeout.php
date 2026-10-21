<?php
$pdo = new PDO('sqlite::memory:');
echo $pdo->query('PRAGMA busy_timeout')->fetchColumn(), '|';
echo $pdo->setAttribute(PDO::ATTR_TIMEOUT, 2) ? 'set|' : 'wrong|';
echo $pdo->query('PRAGMA busy_timeout')->fetchColumn(), '|';
$configured = new PDO('sqlite::memory:', null, null, [PDO::ATTR_TIMEOUT => 1]);
echo $configured->query('PRAGMA busy_timeout')->fetchColumn(), '|';
try {
    $configured->getAttribute(PDO::ATTR_TIMEOUT);
    echo 'wrong';
} catch (PDOException $e) {
    echo $e->getCode();
}
