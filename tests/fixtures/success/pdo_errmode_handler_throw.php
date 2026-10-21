<?php
$pdo = new PDO('sqlite::memory:');
$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_WARNING);
set_error_handler(function ($level, $message) {
    throw new Exception('stop');
});
try {
    $pdo->exec('INSERT INTO absent VALUES (1)');
    echo 'continued|';
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
restore_error_handler();
echo $pdo->errorCode();
