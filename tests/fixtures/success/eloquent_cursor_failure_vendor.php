<?php
require 'research/laravel-api/vendor/autoload.php';

$capsule = new Illuminate\Database\Capsule\Manager;
$capsule->addConnection(['driver' => 'sqlite', 'database' => ':memory:', 'prefix' => '']);
$db = $capsule->getConnection();
$cursor = $db->cursor('SELECT * FROM missing_table');
echo 'created|';
try {
    foreach ($cursor as $row) {
        echo 'unexpected|';
    }
} catch (Illuminate\Database\QueryException $e) {
    echo get_class($e->getPrevious()), '|', strpos($e->getMessage(), '(Connection:') !== false ? 'formatted' : 'unformatted';
}
