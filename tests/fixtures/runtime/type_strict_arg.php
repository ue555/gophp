<?php
declare(strict_types=1);
function typedLocal(int $x): int { return $x; }
echo "before\n";
typedLocal("7");
echo "after\n";
