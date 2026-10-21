<?php
class CompoundFailureState {
    public int $number = 5;
    public float $decimal = 5.5;
    public string $text = 'a';
}
$box = new CompoundFailureState;
try { $box->number /= 0; } catch (DivisionByZeroError $e) { echo 'division:', $box->number, '|'; }
try { $box->number %= 0; } catch (DivisionByZeroError $e) { echo 'modulo:', $box->number, '|'; }
try { $box->number += []; } catch (TypeError $e) { echo 'type:', $box->number, '|'; }
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $box->decimal += '2x'; } catch (Exception $e) { echo 'warning:', $box->decimal; }
restore_error_handler();
echo '|concat:', ($box->text .= 'b'), ':', $box->text;
$plain = 5.5;
$items = [5.5];
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $plain %= 2; } catch (Exception $e) {}
try { $items[0] %= 2; } catch (Exception $e) {}
restore_error_handler();
echo '|plain:', $plain, '|item:', $items[0];
