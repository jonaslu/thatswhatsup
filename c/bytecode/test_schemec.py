import subprocess


def test_assert(expected, actual):
    if not expected == actual:
        raise Exception("Fail: expected " + expected + " actual " + actual)


def get_schemec_output():
    completed_process = subprocess.run(
        ["python", "schemec.py"], stdout=subprocess.PIPE, timeout=10, encoding="utf-8")

    return completed_process.stdout


test_assert(get_schemec_output(), '7')
