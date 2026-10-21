<?php
$pdo = new PDO('sqlite::memory:');
try {
    $pdo->query('SELECT * FROM missing_table');
    echo 'unexpected';
} catch (PDOException $e) {
    echo get_class($e), '|', $e instanceof RuntimeException ? 'runtime' : 'other';
}
