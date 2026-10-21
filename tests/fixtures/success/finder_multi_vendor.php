<?php
require 'research/laravel-api/vendor/autoload.php';
$finder = \Symfony\Component\Finder\Finder::create()->files()
    ->in('tests/fixtures/support/spl_walk')
    ->in('tests/fixtures/support/spl_walk/sub');
$count = 0;
foreach ($finder as $file) {
    ++$count;
}
echo $count;
