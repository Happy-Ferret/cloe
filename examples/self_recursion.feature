Feature: Self recursion
  Scenario: Define a recursive function
    Given a file named "main.cloe" with:
    """
    (def (factorial n)
         (if (= n 0) 1 (* n (factorial (- n 1)))))

    (write (factorial 5))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    120
    """

  Scenario: Define a recursive function in a function definition
    Given a file named "main.cloe" with:
    """
    (def (createFactorial)
         (def (factorial n)
              (if (= n 0) 1 (* n (factorial (- n 1)))))
         factorial)

    (write ((createFactorial) 5))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    120
    """
