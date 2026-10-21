<?php
echo "before\n";
class BadFailure extends Exception {
    public function getTrace() { return []; }
}
