<?php
$pdo = new PDO('sqlite::memory:');
echo $pdo->getAttribute(PDO::ATTR_DEFAULT_FETCH_MODE), '|';
$pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);
$first = $pdo->prepare('SELECT 1 AS Number');
$pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_NUM);
$second = $pdo->prepare('SELECT 2 AS Number');
$first->execute();
$second->execute();
$a = $first->fetch();
$b = $second->fetch();
echo isset($a['Number']) ? 'assoc|' : 'wrong|';
echo isset($b[0]) && !isset($b['Number']) ? 'num|' : 'wrong|';
echo $pdo->getAttribute(PDO::ATTR_DEFAULT_FETCH_MODE), '|';
echo $pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, 99) ? 'stored|' : 'wrong|';
echo $pdo->getAttribute(PDO::ATTR_DEFAULT_FETCH_MODE), '|';
try {
    $pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, 0);
    echo 'wrong';
} catch (ValueError $e) {
    echo 'rejected';
}
