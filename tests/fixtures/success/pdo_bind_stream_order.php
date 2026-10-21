<?php
$db = new PDO('sqlite::memory:');
$stream = fopen('php://temp', 'w+');
fwrite($stream, 'AB');
fseek($stream, 0);
$stmt = $db->prepare('SELECT hex(?) AS first, hex(?) AS second');
$stmt->bindValue(2, $stream, PDO::PARAM_LOB);
$stmt->bindValue(1, $stream, PDO::PARAM_LOB);
$stmt->execute();
$row = $stmt->fetch(PDO::FETCH_ASSOC);
echo $row['first'], '|', $row['second'];
fclose($stream);
