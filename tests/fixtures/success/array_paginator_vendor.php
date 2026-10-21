<?php
require 'research/laravel-api/vendor/autoload.php';
$paginator = new \ApiPlatform\State\Pagination\ArrayPaginator([10, 20, 30], 1, 1);
foreach ($paginator as $key => $value) echo $key, ':', $value, ';';
echo '|', count($paginator), '|', $paginator->getTotalItems();
