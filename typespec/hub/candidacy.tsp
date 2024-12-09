import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddHubCandidacyCommentRequest {
    candidacy_id: string;

    @maxLength(2048)
    comment: string;
}

model MyCandidacy {
    candidacy_id: string;
    company_name: string;
    company_domain: string;
    opening_id: string;
    opening_title: string;
    opening_description: string;
    candidacy_state: CandidacyState;
}

@route("/hub/add-candidacy-comment")
interface AddHubCandidacyComment {
    @tag("Candidacies")
    @post
    addComment(@body request: AddHubCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/get-my-candidacies")
interface GetMyCandidacies {
    @tag("Candidacies")
    @get
    getMyCandidacies(): {
        @statusCode statusCode: 200;
    };
}
