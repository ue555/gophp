<?php
trait AliasFirst { public function first() { return 1; } }
trait AliasSecond { public function second() { return 2; } }
echo "before\n";
if (true) {
    class AliasCollision {
        use AliasFirst, AliasSecond {
            AliasFirst::first as same;
            AliasSecond::second as same;
        }
    }
}
echo "after\n";
