<?php
$pdo = new PDO('sqlite::memory:');
echo $pdo->getAttribute(PDO::ATTR_DRIVER_NAME), '|', $pdo->getAttribute(PDO::ATTR_CASE), '|';
$pdo->setAttribute(PDO::ATTR_CASE, PDO::CASE_UPPER);
$pdo->setAttribute(PDO::ATTR_STRINGIFY_FETCHES, true);
$pdo->setAttribute(PDO::ATTR_ORACLE_NULLS, PDO::NULL_EMPTY_STRING);
$row = $pdo->query("SELECT 7 AS MiXeD, 2.5 AS FlOat, NULL AS Missing, '' AS Empty")->fetch(PDO::FETCH_ASSOC);
echo gettype($row['MIXED']), '|', $row['MIXED'], '|', gettype($row['FLOAT']), '|', $row['FLOAT'], '|';
echo $row['EMPTY'] === null ? 'empty-null|' : 'wrong|';
echo $pdo->getAttribute(PDO::ATTR_STRINGIFY_FETCHES) ? 'stringify|' : 'wrong|';
$pdo->setAttribute(PDO::ATTR_ORACLE_NULLS, PDO::NULL_TO_STRING);
$row = $pdo->query('SELECT NULL AS Missing')->fetch(PDO::FETCH_ASSOC);
echo $row['MISSING'] === '' ? 'null-string|' : 'wrong|';
echo $pdo->setAttribute(PDO::ATTR_EMULATE_PREPARES, false) === false ? 'unsupported|' : 'wrong|';

$configured = new PDO('sqlite::memory:', null, null, [PDO::ATTR_CASE => PDO::CASE_LOWER]);
$row = $configured->query('SELECT 1 AS Mixed')->fetch(PDO::FETCH_ASSOC);
echo $configured->getAttribute(PDO::ATTR_CASE), '|', $row['mixed'];
try {
    $pdo->getAttribute(PDO::ATTR_EMULATE_PREPARES);
    echo '|wrong';
} catch (PDOException $e) {
    echo '|', $e->getCode(), '|', $e->errorInfo[1], '|', $e->errorInfo[2] === null ? 'null' : 'wrong';
}

$pdo->setAttribute(PDO::ATTR_CASE, PDO::CASE_UPPER);
$snapshot = $pdo->query('SELECT 3 AS SwItCh');
$pdo->setAttribute(PDO::ATTR_CASE, PDO::CASE_LOWER);
$row = $snapshot->fetch(PDO::FETCH_ASSOC);
echo '|', isset($row['SWITCH']) ? 'snapshot' : 'wrong';

$prepared = $pdo->prepare('SELECT 4 AS ReAdY');
$pdo->setAttribute(PDO::ATTR_CASE, PDO::CASE_UPPER);
$prepared->execute();
$row = $prepared->fetch(PDO::FETCH_ASSOC);
echo '|', isset($row['READY']) ? 'execute' : 'wrong';
