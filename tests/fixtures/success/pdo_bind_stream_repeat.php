<?php
$db = new PDO('sqlite::memory:');
$stream = fopen('php://temp', 'w+');
fwrite($stream, "Q\0R");
fseek($stream, 0);
$stmt = $db->prepare('SELECT hex(?)');
$stmt->bindValue(1, $stream, PDO::PARAM_LOB);
$stmt->execute();
echo $stmt->fetchColumn(), '|';
$stmt->execute();
echo $stmt->fetchColumn();
fclose($stream);
