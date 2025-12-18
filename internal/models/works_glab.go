package models

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Group struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	//Visibility  string `json:"visibility"`
	// ShareWithGroupLock                bool   `json:"share_with_group_lock"`
	// Require_two_factor_authentication bool   `json:"require_two_factor_authentication"`
	// Two_factor_grace_period           int    `json:"two_factor_grace_period"`
	Project_creation_level string `json:"project_creation_level"`
	ParentId               int    `json:"parent_id"`
	CreatedAt              string `json:"created_at"`
	/*"auto_devops_enabled": null,
	  "subgroup_creation_level": "owner",
	  "emails_disabled": null,
	  "emails_enabled": null,
	  "mentions_disabled": null,
	  "lfs_enabled": true,
	  "default_branch": null,
	  "default_branch_protection": 2,
	  "default_branch_protection_defaults": {
	    "allowed_to_push": [
	        {
	            "access_level": 40
	        }
	    ],
	    "allow_force_push": false,
	    "allowed_to_merge": [
	        {
	            "access_level": 40
	        }
	    ]
	  },
	  "avatar_url": "http://gitlab.example.com/uploads/group/avatar/1/foo.jpg",
	  "web_url": "http://gitlab.example.com/groups/foo-bar",
	  "request_access_enabled": false,
	  "repository_storage": "default",
	  "full_name": "Foobar Group",
	  "full_path": "foo-bar",
	  "file_template_project_id": 1,
	*/

}

func (m *WorkModel) Sync(w http.ResponseWriter, r *http.Request, ip string, token string) ([]Group, error) {
	client := &http.Client{}
	pages := 1
	//var gr []int
	var allGroups []Group
	// данных может быть много, поэтому запрашиваем пока не получим пусто
	for {
		req, err := http.NewRequest(
			http.MethodGet, ip+"/api/v4/groups?per_page=100&page="+strconv.Itoa(pages), nil,
		)

		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			// return fmt.Errorf("error sending a request %v", err)
			return nil, err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			//h.Logger.Sugar().Infof("error reading a request %v", err)
			return nil, err
		}

		if len(body) < 3 {
			break
		}
		//fmt.Println(string(body), len(body))

		var grs []Group
		err = json.Unmarshal(body, &grs)
		if err != nil {
			//fmt.Printf("Ошибка при декодировании %v\n", err)
			return nil, err
		}

		for _, item := range grs {
			if item.ParentId == 0 {
				//			gr = append(gr, item.Id)
				allGroups = append(allGroups, item)
				//	fmt.Println("\n", item.Name, item.Id, item.Path, item.Project_creation_level)
			}
		}
		//fmt.Println("pages = ", pages)
		pages++
	}

	return allGroups, nil
}
