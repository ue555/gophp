<?php
$pdo = new PDO('sqlite::memory:');
try {
    echo $pdo->setAttribute(PDO::ATTR_ORACLE_NULLS, 9) ? 'accepted|' : 'false|';
} catch (ValueError $e) {
    echo 'rejected|';
}
echo $pdo->getAttribute(PDO::ATTR_ORACLE_NULLS);
