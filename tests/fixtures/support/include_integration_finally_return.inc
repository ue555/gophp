<?php
function includedDecision() {
    try {
        throw new Exception('discarded');
    } finally {
        return 9;
    }
}
return includedDecision();
