<?php
$pdo = new PDO('sqlite::memory:');
echo $pdo->getAttribute(PDO::ATTR_ERRMODE), '|';
$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_SILENT);
echo $pdo->getAttribute(PDO::ATTR_ERRMODE), '|';
echo $pdo->exec('INSERT INTO absent VALUES (1)') === false ? 'silent|' : 'wrong|';
echo $pdo->errorCode(), '|';

set_error_handler(function ($level, $message) {
    echo $level === E_WARNING ? 'warning|' : 'wrong-level|';
    return true;
});
$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_WARNING);
echo $pdo->exec('INSERT INTO absent VALUES (1)') === false ? 'warn-false|' : 'wrong|';
restore_error_handler();
echo $pdo->errorCode(), '|';

try {
    $pdo->setAttribute(PDO::ATTR_ERRMODE, 99);
} catch (ValueError $e) {
    echo 'invalid|';
}
echo $pdo->getAttribute(PDO::ATTR_ERRMODE), '|';

$configured = new PDO('sqlite::memory:', null, null, [PDO::ATTR_ERRMODE => PDO::ERRMODE_SILENT]);
echo $configured->getAttribute(PDO::ATTR_ERRMODE);
