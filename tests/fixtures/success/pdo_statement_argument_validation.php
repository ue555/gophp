<?php
$db = new PDO('sqlite::memory:');
$stmt = $db->prepare('SELECT ?');
$stmt->bindValue(1, 7, PDO::PARAM_INT);
try {
    $stmt->bindValue(0, 9);
} catch (ValueError $e) {
    echo 'zero|';
}
try {
    $stmt->bindValue('', 9);
} catch (ValueError $e) {
    echo 'empty|';
}
try {
    $stmt->bindValue([], 9);
} catch (TypeError $e) {
    echo 'param-type|';
}
try {
    $stmt->execute('wrong');
} catch (TypeError $e) {
    echo 'execute-type|';
}
$stmt->execute();
echo $stmt->fetchColumn(), '|';
$stmt->execute(params: [12]);
echo $stmt->fetchColumn();
