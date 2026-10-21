<?php
$db = new PDO('sqlite::memory:');
$stmt = $db->prepare('SELECT :known');
$stmt->bindValue(':missing', 'x');
try {
    $stmt->execute();
    echo 'unexpected|';
} catch (PDOException $e) {
    $info = $stmt->errorInfo();
    echo $e->getCode(), '|', $stmt->errorCode(), '|', $info[1], '|', $info[2], '|';
}
$stmt->execute(['known' => 'ready']);
echo $stmt->errorCode(), '|', $stmt->fetchColumn();
