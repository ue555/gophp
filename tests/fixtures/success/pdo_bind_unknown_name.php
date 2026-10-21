<?php
$db = new PDO('sqlite::memory:');
$stmt = $db->prepare('SELECT :known');
echo $stmt->bindValue(':missing', 'x') ? 'registered|' : 'rejected|';
try {
    $stmt->execute();
    echo 'executed|';
} catch (PDOException $e) {
    echo 'execute-error|';
}
$stmt->execute(['known' => 'ready']);
echo $stmt->fetchColumn();
