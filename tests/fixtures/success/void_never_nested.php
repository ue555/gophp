<?php
function outer(): void {
    function inner(): int { return 4; }
    $c = function (): int { return 5; };
    echo inner(), '|', $c(), '|';
}
outer();
$arrow = fn(): never => throw new Exception('arrow');
try { $arrow(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
