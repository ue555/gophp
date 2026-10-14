<?php
trait FirstSeparate { public function value() { return 1; } }
trait SecondSeparate { public function value() { return 2; } }
class SeparateCollision {
    use FirstSeparate;
    use SecondSeparate;
}
echo (new SeparateCollision)->value();
