For each of the following 1 API endpoints:

API 1: GET /graphql (Save as api_combined_0_0.txt)

API 1 is found in "/examples/context/main.go"

For each API endpoint, get its input, output, sample input, and sample output. As well as the relevant code snippets to fully understand it.

For each API endpoint, do this:

1. Find the implementation of the API endpoint by searching across all files.

2. Extract and include ALL relevant code snippets:
   - The route definition/declaration
   - The controller/handler function complete implementation
   - Any middleware functions used
   - Any helper functions or services called
   - Database models or queries used

For each code snippet, include the file path and line numbers where the code is found.

3. Based on analyzing these code snippets, determine:

   - Input Format: 
     * URL parameters
     * Query parameters
     * Request body fields and structure
     * Required vs optional fields
     * Data types for each field

   - Output Format:
     * Response structure
     * All possible status codes
     * Data types for each field
     * Error response formats

   - Sample Input: A realistic curl command that would work with this API

   - Sample Output: The expected JSON response for the sample input

The whole file shall be in markdown format. 

The method and path of the API endpoint should be used as the heading 1.

The "Implementation", "Input Format", "Output Format", "Sample Input", and "Sample Output" sections should be clearly marked as heading 2.

IMPORTANT: Save your analysis of EACH API endpoint as a SEPARATE file in the /elasticdash/ directory using the exact filename specified next to each API in the list above.

I need separate, complete files for each API endpoint, with each file containing all the sections above.