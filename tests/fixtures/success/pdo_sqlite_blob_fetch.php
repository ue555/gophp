<?php
$pdo = new PDO('sqlite::memory:');
$row = $pdo->query("SELECT X'410042' AS Payload, X'' AS EmptyValue, NULL AS Missing")->fetch(PDO::FETCH_ASSOC);
echo gettype($row['Payload']), '|', bin2hex($row['Payload']), '|';
echo gettype($row['EmptyValue']), '|', strlen($row['EmptyValue']), '|';
echo gettype($row['Missing']);
