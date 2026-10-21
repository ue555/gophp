<?php
$pdo = new PDO('sqlite::memory:');
$statement = $pdo->query("SELECT 1 AS First, 2 AS Second UNION ALL SELECT 3, 4");
echo $statement->fetchColumn(1), '|', $statement->fetchColumn(0), '|';
echo $statement->fetchColumn() === false ? 'end|' : 'wrong|';
$statement = $pdo->query('SELECT 5 AS First');
try {
    $statement->fetchColumn(-1);
    echo 'wrong|';
} catch (ValueError $e) {
    echo 'negative|';
}
$statement = $pdo->query('SELECT 6 AS First');
try {
    $statement->fetchColumn(1);
    echo 'wrong|';
} catch (ValueError $e) {
    echo 'missing|';
}
$pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);
$statement = $pdo->query('SELECT 7 AS Named');
$statement->setFetchMode(PDO::FETCH_NUM);
$statement->setFetchMode(PDO::FETCH_DEFAULT);
$row = $statement->fetch();
echo isset($row[0]) ? 'both' : 'assoc';
