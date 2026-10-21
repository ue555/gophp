<?php
function importedCallable(int $n): int { return $n + 10; }
class ImportedCallable {
    public static function run(int $n): int { return $n + 20; }
}
